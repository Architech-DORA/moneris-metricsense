package services

import (
	"bytes"
	"text/template"

	"github.com/apache/incubator-devlake/core/errors"
)

type MetricsQuery struct{}

type DeploymentFrequency struct {
	Summary interface{} `json:"summary"`
}

type Metrics struct {
	DeploymentFrequency map[string]interface{} `json:"deploymentFrequency"`
}

const DEPLOYMENT_FREQUENCY_QUERY = `
-- Metric 1: Deployment Frequency
with last_few_calendar_months as(
-- construct the last few calendar months within the selected time period in the top-right corner
	SELECT CAST((SYSDATE()-INTERVAL (H+T+U) DAY) AS date) day
	FROM ( SELECT 0 H
			UNION ALL SELECT 100 UNION ALL SELECT 200 UNION ALL SELECT 300
		) H CROSS JOIN ( SELECT 0 T
			UNION ALL SELECT  10 UNION ALL SELECT  20 UNION ALL SELECT  30
			UNION ALL SELECT  40 UNION ALL SELECT  50 UNION ALL SELECT  60
			UNION ALL SELECT  70 UNION ALL SELECT  80 UNION ALL SELECT  90
		) T CROSS JOIN ( SELECT 0 U
			UNION ALL SELECT   1 UNION ALL SELECT   2 UNION ALL SELECT   3
			UNION ALL SELECT   4 UNION ALL SELECT   5 UNION ALL SELECT   6
			UNION ALL SELECT   7 UNION ALL SELECT   8 UNION ALL SELECT   9
		) U
	WHERE
		(SYSDATE()-INTERVAL (H+T+U) DAY) > FROM_UNIXTIME({{.dateFrom}})
),

_production_deployment_days as(
-- When deploying multiple commits in one pipeline, GitLab and BitBucket may generate more than one deployment. However, DevLake consider these deployments as ONE production deployment and use the last one's finished_date as the finished date.
	SELECT
		cdc.cicd_deployment_id as deployment_id,
		max(DATE(cdc.finished_date)) as day
	FROM cicd_deployment_commits cdc
	JOIN project_mapping pm on cdc.cicd_scope_id = pm.row_id and pm.table = 'cicd_scopes'
	WHERE
		pm.project_name in ({{.projects}})
		and cdc.result = 'SUCCESS'
		and cdc.environment = 'PRODUCTION'
	GROUP BY 1
),

_days_weeks_deploy as(
-- calculate the number of deployment days every week
	SELECT
			date(DATE_ADD(last_few_calendar_months.day, INTERVAL -WEEKDAY(last_few_calendar_months.day) DAY)) as week,
			MAX(if(_production_deployment_days.day is not null, 1, 0)) as weeks_deployed,
			COUNT(distinct _production_deployment_days.day) as days_deployed
	FROM 
		last_few_calendar_months
		LEFT JOIN _production_deployment_days ON _production_deployment_days.day = last_few_calendar_months.day
	GROUP BY week
	),

_monthly_deploy as(
-- calculate the number of deployment days every month
	SELECT
			date(DATE_ADD(last_few_calendar_months.day, INTERVAL -DAY(last_few_calendar_months.day)+1 DAY)) as month,
			MAX(if(_production_deployment_days.day is not null, 1, 0)) as months_deployed
	FROM 
		last_few_calendar_months
		LEFT JOIN _production_deployment_days ON _production_deployment_days.day = last_few_calendar_months.day
	GROUP BY month
	),

_median_number_of_deployment_days_per_week_ranks as(
	SELECT *, percent_rank() over(order by days_deployed) as ranks
	FROM _days_weeks_deploy
),

_median_number_of_deployment_days_per_week as(
	SELECT max(days_deployed) as median_number_of_deployment_days_per_week
	FROM _median_number_of_deployment_days_per_week_ranks
	WHERE ranks <= 0.5
),

_median_number_of_deployment_days_per_month_ranks as(
	SELECT *, percent_rank() over(order by months_deployed) as ranks
	FROM _monthly_deploy
),

_median_number_of_deployment_days_per_month as(
	SELECT max(months_deployed) as median_number_of_deployment_days_per_month
	FROM _median_number_of_deployment_days_per_month_ranks
	WHERE ranks <= 0.5
)

SELECT 
	CASE  
		WHEN median_number_of_deployment_days_per_week >= 3 THEN 'On-demand'
		WHEN median_number_of_deployment_days_per_week >= 1 THEN 'Between once per week and once per month'
		WHEN median_number_of_deployment_days_per_month >= 1 THEN 'Between once per month and once every 6 months'
		ELSE 'Fewer than once per six months' END AS 'deploymentFrequency'
FROM _median_number_of_deployment_days_per_week, _median_number_of_deployment_days_per_month
`

const DEPLOYMENT_FREQUENCY_CHART_DATA_QUERY = `
-- Metric 1: Number of deployments per month
with _deployments as(
-- When deploying multiple commits in one pipeline, GitLab and BitBucket may generate more than one deployment. However, DevLake consider these deployments as ONE production deployment and use the last one's finished_date as the finished date.
	SELECT 
		date_format(deployment_finished_date,'%y/%m') as month,
		count(cicd_deployment_id) as deployment_count
	FROM (
		SELECT
			cdc.cicd_deployment_id,
			max(cdc.finished_date) as deployment_finished_date
		FROM cicd_deployment_commits cdc
		JOIN project_mapping pm on cdc.cicd_scope_id = pm.row_id and pm.table = 'cicd_scopes'
		WHERE
			pm.project_name in ({{.projects}})
			and cdc.result = 'SUCCESS'
			and cdc.environment = 'PRODUCTION'
		GROUP BY 1
		HAVING max(cdc.finished_date BETWEEN FROM_UNIXTIME({{.dateFrom}}) AND FROM_UNIXTIME({{.dateTo}}))
	) _production_deployments
	GROUP BY 1
),

_calendar_months as(
-- construct the calendar months of last 6 months
	SELECT date_format(CAST((SYSDATE()-INTERVAL (month_index) MONTH) AS date), '%y/%m') as month
	FROM ( SELECT 0 month_index
			UNION ALL SELECT   1  UNION ALL SELECT   2 UNION ALL SELECT   3
			UNION ALL SELECT   4  UNION ALL SELECT   5 UNION ALL SELECT   6
			UNION ALL SELECT   7  UNION ALL SELECT   8 UNION ALL SELECT   9
			UNION ALL SELECT   10 UNION ALL SELECT  11
		) month_index
	WHERE (SYSDATE()-INTERVAL (month_index) MONTH) > SYSDATE()-INTERVAL 12 MONTH	
)

SELECT 
	cm.month, 
	case when d.deployment_count is null then 0 else d.deployment_count end as dataCount
FROM 
	_calendar_months cm
	left join _deployments d on cm.month = d.month
ORDER BY 1
`

const MEDIAN_LEAD_TIME_FOR_CHANGES = `
-- Metric 2: median lead time for changes
with _pr_stats as (
-- get the cycle time of PRs deployed by the deployments finished in the selected period
	SELECT
		distinct pr.id,
		ppm.pr_cycle_time
	FROM
		pull_requests pr 
		join project_pr_metrics ppm on ppm.id = pr.id
		join project_mapping pm on pr.base_repo_id = pm.row_id and pm.table = 'repos'
		join cicd_deployment_commits cdc on ppm.deployment_commit_id = cdc.id
	WHERE
	  pm.project_name in ({{.projects}}) 
		and pr.merged_date is not null
		and ppm.pr_cycle_time is not null
		and cdc.finished_date BETWEEN FROM_UNIXTIME({{.dateFrom}}) AND FROM_UNIXTIME({{.dateTo}})
),

_median_change_lead_time_ranks as(
	SELECT *, percent_rank() over(order by pr_cycle_time) as ranks
	FROM _pr_stats
),

_median_change_lead_time as(
-- use median PR cycle time as the median change lead time
	SELECT max(pr_cycle_time) as median_change_lead_time
	FROM _median_change_lead_time_ranks
	WHERE ranks <= 0.5
)
SELECT 
	CASE
		WHEN median_change_lead_time < 60 then "Less than one hour"
		WHEN median_change_lead_time < 7 * 24 * 60 then "Less than one week"
		WHEN median_change_lead_time < 180 * 24 * 60 then "Between one week and six months"
		ELSE "More than six months"
		END as leadTimeForChanges
FROM _median_change_lead_time
;
`

const MEDIAN_LEAD_TIME_FOR_CHANGES_CHART_QUERY = `
	-- Metric 2: median change lead time per month
	with _pr_stats as (
	-- get the cycle time of PRs deployed by the deployments finished each month
		SELECT
			distinct pr.id,
			date_format(cdc.finished_date,'%y/%m') as month,
			ppm.pr_cycle_time
		FROM
			pull_requests pr
			join project_pr_metrics ppm on ppm.id = pr.id
			join project_mapping pm on pr.base_repo_id = pm.row_id and pm.table = 'repos'
			join cicd_deployment_commits cdc on ppm.deployment_commit_id = cdc.id
		WHERE
			pm.project_name in ({{.projects}}) 
			and pr.merged_date is not null
			and ppm.pr_cycle_time is not null
			and cdc.finished_date BETWEEN FROM_UNIXTIME({{.dateFrom}}) AND FROM_UNIXTIME({{.dateTo}})
	),
	
	_find_median_clt_each_month_ranks as(
		SELECT *, percent_rank() over(PARTITION BY month order by pr_cycle_time) as ranks
		FROM _pr_stats
	),
	
	_clt as(
		SELECT month, max(pr_cycle_time) as median_change_lead_time
		FROM _find_median_clt_each_month_ranks
		WHERE ranks <= 0.5
		group by month
	),
	
	_calendar_months as(
	-- to	deal with the month with no incidents
		SELECT date_format(CAST((SYSDATE()-INTERVAL (month_index) MONTH) AS date), '%y/%m') as month
		FROM ( SELECT 0 month_index
				UNION ALL SELECT   1  UNION ALL SELECT   2 UNION ALL SELECT   3
				UNION ALL SELECT   4  UNION ALL SELECT   5 UNION ALL SELECT   6
				UNION ALL SELECT   7  UNION ALL SELECT   8 UNION ALL SELECT   9
				UNION ALL SELECT   10 UNION ALL SELECT  11
			) month_index
		WHERE (SYSDATE()-INTERVAL (month_index) MONTH) > SYSDATE()-INTERVAL 12 MONTH	
	)
	
	SELECT 
		cm.month,
		case 
			when _clt.median_change_lead_time is null then 0 
			else _clt.median_change_lead_time/60 end as medianChangeLeadTimeInHour
	FROM 
		_calendar_months cm
		left join _clt on cm.month = _clt.month
	ORDER BY 1
`

const MEDIAN_TIME_TO_RESTORE = `
-- Metric 3: Median time to restore service 
WITH _incidents as (
-- get the incidents created within the selected time period in the top-right corner
	SELECT
	  distinct i.id,
		cast(lead_time_minutes as signed) as lead_time_minutes
	FROM
		issues i
	  join board_issues bi on i.id = bi.issue_id
	  join boards b on bi.board_id = b.id
	  join project_mapping pm on b.id = pm.row_id and pm.table = 'boards'
	WHERE
	  pm.project_name in ({{.projects}})
		and i.type = 'INCIDENT'
		and i.created_date BETWEEN FROM_UNIXTIME({{.dateFrom}}) AND FROM_UNIXTIME({{.dateTo}})
),

_median_mttr_ranks as(
	SELECT *, percent_rank() over(order by lead_time_minutes) as ranks
	FROM _incidents
),

_median_mttr as(
	SELECT max(lead_time_minutes) as median_time_to_resolve
	FROM _median_mttr_ranks
	WHERE ranks <= 0.5
)
SELECT 
	case
		WHEN median_time_to_resolve < 60  then "Less than one hour"
		WHEN median_time_to_resolve < 24 * 60 then "Less than one Day"
		WHEN median_time_to_resolve < 7 * 24 * 60  then "Between one day and one week"
		ELSE "More than one week"
		END as "timeToRestoreService"
FROM 
	_median_mttr;`

const MEDIAN_TIME_TO_RESTORE_SERVICE_CHART_QUERY = `
-- Metric 3: median time to restore service - MTTR
with _incidents as (
-- get the number of incidents created each month
	SELECT
	  distinct i.id,
		date_format(i.created_date,'%y/%m') as month,
		cast(lead_time_minutes as signed) as lead_time_minutes
	FROM
		issues i
	  join board_issues bi on i.id = bi.issue_id
	  join boards b on bi.board_id = b.id
	  join project_mapping pm on b.id = pm.row_id and pm.table = 'boards'
	WHERE
	  pm.project_name in ({{.projects}})
		and i.type = 'INCIDENT'
		and i.lead_time_minutes is not null
),

_find_median_mttr_each_month_ranks as(
	SELECT *, percent_rank() over(PARTITION BY month order by lead_time_minutes) as ranks
	FROM _incidents
),

_mttr as(
	SELECT month, max(lead_time_minutes) as median_time_to_resolve
	FROM _find_median_mttr_each_month_ranks
	WHERE ranks <= 0.5
	GROUP BY month
),

_calendar_months as(
-- deal with the month with no incidents
	SELECT date_format(CAST((SYSDATE()-INTERVAL (month_index) MONTH) AS date), '%y/%m') as month
	FROM ( SELECT 0 month_index
			UNION ALL SELECT   1  UNION ALL SELECT   2 UNION ALL SELECT   3
			UNION ALL SELECT   4  UNION ALL SELECT   5 UNION ALL SELECT   6
			UNION ALL SELECT   7  UNION ALL SELECT   8 UNION ALL SELECT   9
			UNION ALL SELECT   10 UNION ALL SELECT  11
		) month_index
	WHERE (SYSDATE()-INTERVAL (month_index) MONTH) > SYSDATE()-INTERVAL 12 MONTH	
)

SELECT 
	cm.month,
	case 
		when m.median_time_to_resolve is null then 0 
		else m.median_time_to_resolve/60 end as medianTimeToResolveInHour
FROM 
	_calendar_months cm
	left join _mttr m on cm.month = m.month
ORDER BY 1`

const CHANGE_FAILURE_RATE = `
-- Metric 4: change failure rate
WITH _deployments as (
-- When deploying multiple commits in one pipeline, GitLab and BitBucket may generate more than one deployment. However, DevLake consider these deployments as ONE production deployment and use the last one's finished_date as the finished date.
	SELECT
		cdc.cicd_deployment_id as deployment_id,
		max(cdc.finished_date) as deployment_finished_date
	FROM 
		cicd_deployment_commits cdc
		JOIN project_mapping pm on cdc.cicd_scope_id = pm.row_id and pm.table = 'cicd_scopes'
	WHERE
		pm.project_name in ({{.projects}})
		and cdc.result = 'SUCCESS'
		and cdc.environment = 'PRODUCTION'
	GROUP BY 1
	HAVING max(cdc.finished_date BETWEEN FROM_UNIXTIME({{.dateFrom}}) AND FROM_UNIXTIME({{.dateTo}}))
),

_failure_caused_by_deployments as (
-- calculate the number of incidents caused by each deployment
	SELECT
		d.deployment_id,
		d.deployment_finished_date,
		count(distinct case when i.type = 'INCIDENT' then d.deployment_id else null end) as has_incident
	FROM
		_deployments d
		left join project_issue_metrics pim on d.deployment_id = pim.deployment_id
		left join issues i on pim.id = i.id
	GROUP BY 1,2
),

_change_failure_rate as (
	SELECT 
		case 
			when count(deployment_id) is null then null
			else sum(has_incident)/count(deployment_id) end as change_failure_rate
	FROM
		_failure_caused_by_deployments
)

SELECT
	case  
		when change_failure_rate <= .15 then "0-15%"
		when change_failure_rate <= .20 then "16%-20%"
		when change_failure_rate <= .30 then "21%-30%"
		else "> 30%" 
	end as "changeFailureRate"
FROM 
	_change_failure_rate
;`

const CHANGE_FAILURE_RATE_CHART_QUERY = `
-- Metric 4: change failure rate per month
with _deployments as (
-- When deploying multiple commits in one pipeline, GitLab and BitBucket may generate more than one deployment. However, DevLake consider these deployments as ONE production deployment and use the last one's finished_date as the finished date.
	SELECT
		cdc.cicd_deployment_id as deployment_id,
		max(cdc.finished_date) as deployment_finished_date
	FROM 
		cicd_deployment_commits cdc
		JOIN project_mapping pm on cdc.cicd_scope_id = pm.row_id and pm.table = 'cicd_scopes'
	WHERE
		pm.project_name in ({{.projects}})
		and cdc.result = 'SUCCESS'
		and cdc.environment = 'PRODUCTION'
	GROUP BY 1
	HAVING max(cdc.finished_date BETWEEN FROM_UNIXTIME({{.dateFrom}}) AND FROM_UNIXTIME({{.dateTo}}))
),

_failure_caused_by_deployments as (
-- calculate the number of incidents caused by each deployment
	SELECT
		d.deployment_id,
		d.deployment_finished_date,
		count(distinct case when i.type = 'INCIDENT' then d.deployment_id else null end) as has_incident
	FROM
		_deployments d
		left join project_issue_metrics pim on d.deployment_id = pim.deployment_id
		left join issues i on pim.id = i.id
	GROUP BY 1,2
),

_change_failure_rate_for_each_month as (
	SELECT 
		date_format(deployment_finished_date,'%y/%m') as month,
		case 
			when count(deployment_id) is null then null
			else sum(has_incident)/count(deployment_id) end as change_failure_rate
	FROM
		_failure_caused_by_deployments
	GROUP BY 1
),

_calendar_months as(
-- deal with the month with no incidents
	SELECT date_format(CAST((SYSDATE()-INTERVAL (month_index) MONTH) AS date), '%y/%m') as month
	FROM ( SELECT 0 month_index
			UNION ALL SELECT   1  UNION ALL SELECT   2 UNION ALL SELECT   3
			UNION ALL SELECT   4  UNION ALL SELECT   5 UNION ALL SELECT   6
			UNION ALL SELECT   7  UNION ALL SELECT   8 UNION ALL SELECT   9
			UNION ALL SELECT   10 UNION ALL SELECT  11
		) month_index
	WHERE (SYSDATE()-INTERVAL (month_index) MONTH) > SYSDATE()-INTERVAL 12 MONTH	
)

SELECT 
	cm.month,
	cfr.change_failure_rate AS changeFailureRate
FROM 
	_calendar_months cm
	left join _change_failure_rate_for_each_month cfr on cm.month = cfr.month
GROUP BY 1,2
ORDER BY 1 `

func formatQueryString(name string, queryStr string, args map[string]interface{}) (string, error) {
	var formattedQuery bytes.Buffer
	tmpl, err := template.New(name).Parse(queryStr)

	if err != nil {
		return "", err
	}

	tmpl.Execute(&formattedQuery, args)

	return formattedQuery.String(), nil
}

func GetMetrics(params map[string]interface{}) (*Metrics, errors.Error) {
	metricsData := &Metrics{}

	deploymentFrequencyQueryString, formatError := formatQueryString("deploymentFrequencyQueryString", DEPLOYMENT_FREQUENCY_QUERY, params)
	if formatError != nil {
		return nil, errors.Default.Wrap(formatError, "Something went wrong!")
	}
	rows, err := db.RawCursor(deploymentFrequencyQueryString)
	for rows.Next() {
		data := make(map[string]interface{})
		db.Fetch(rows, data)
		metricsData.DeploymentFrequency = data
	}
	if err != nil {
		return nil, errors.Default.Wrap(err, "Something went wrong!")
	}

	return metricsData, nil
}

type DeploymentFrequencyChartData struct {
	deploymentFinishedDate string
	dataCount              int
}

func GetDeploymentFrequency(params map[string]interface{}) (*map[string]interface{}, errors.Error) {
	deploymentFrequencyQueryString, formatError := formatQueryString("deploymentFrequencyQueryString", DEPLOYMENT_FREQUENCY_QUERY, params)
	if formatError != nil {
		return nil, errors.Default.Wrap(formatError, "Something went wrong!")
	}

	rows, err := db.RawCursor(deploymentFrequencyQueryString)
	defer rows.Close()
	if err != nil {
		return nil, errors.Default.Wrap(err, "Something went wrong!")
	}

	data := make(map[string]interface{})
	for rows.Next() {
		db.Fetch(rows, data)
	}

	deploymentFrequencyChartQueryString, formatError := formatQueryString("deploymentFrequencyChartQueryString", DEPLOYMENT_FREQUENCY_CHART_DATA_QUERY, params)
	if formatError != nil {
		return nil, errors.Default.Wrap(formatError, "Something went wrong!")
	}

	chartDataRows, err := db.RawCursor(deploymentFrequencyChartQueryString)

	if err != nil {
		return nil, errors.Default.Wrap(err, "Something went wrong!")
	}

	var chartData []interface{}

	for chartDataRows.Next() {
		data := make(map[string]interface{})

		db.Fetch(chartDataRows, data)

		chartData = append(chartData, data)

	}

	data["chartData"] = chartData

	return &data, nil
}

func GetLeadTimeForChanges(params map[string]interface{}) (*map[string]interface{}, errors.Error) {
	leadTimeForChangesQueryString, formatError := formatQueryString("leadTimeForChangesQueryString", MEDIAN_LEAD_TIME_FOR_CHANGES, params)
	if formatError != nil {
		return nil, errors.Default.Wrap(formatError, "Something went wrong!")
	}

	rows, err := db.RawCursor(leadTimeForChangesQueryString)
	if err != nil {
		return nil, errors.Default.Wrap(err, "Something went wrong!")
	}
	data := make(map[string]interface{})
	for rows.Next() {
		db.Fetch(rows, data)
	}

	leadTimeForChangesChartQueryString, formatError := formatQueryString("leadTimeForChangesChartQueryString", MEDIAN_LEAD_TIME_FOR_CHANGES_CHART_QUERY, params)
	if formatError != nil {
		return nil, errors.Default.Wrap(formatError, "Something went wrong!")
	}
	chartDataRows, err := db.RawCursor(leadTimeForChangesChartQueryString)

	if err != nil {
		return nil, errors.Default.Wrap(err, "Something went wrong!")
	}

	var chartData []interface{}

	for chartDataRows.Next() {
		data := make(map[string]interface{})

		db.Fetch(chartDataRows, data)

		chartData = append(chartData, data)

	}

	data["chartData"] = chartData
	return &data, nil
}

func GetMedianTimeToRestore(params map[string]interface{}) (*map[string]interface{}, errors.Error) {
	medianTimeToRestoreServiceQueryString, formatError := formatQueryString("medianTimeToRestoreServiceQueryString", MEDIAN_TIME_TO_RESTORE, params)
	if formatError != nil {
		return nil, errors.Default.Wrap(formatError, "Something went wrong!")
	}
	rows, err := db.RawCursor(medianTimeToRestoreServiceQueryString)

	if err != nil {
		return nil, errors.Default.Wrap(err, "Something went wrong!")
	}
	data := make(map[string]interface{})
	for rows.Next() {
		db.Fetch(rows, data)
	}
	medianTimeToRestoreServiceChartQueryString, formatError := formatQueryString("medianTimeToRestoreServiceChartQueryString", MEDIAN_TIME_TO_RESTORE_SERVICE_CHART_QUERY, params)
	if formatError != nil {
		return nil, errors.Default.Wrap(formatError, "Something went wrong!")
	}
	chartDataRows, err := db.RawCursor(medianTimeToRestoreServiceChartQueryString)

	if err != nil {
		return nil, errors.Default.Wrap(err, "Something went wrong!")
	}

	var chartData []interface{}

	for chartDataRows.Next() {
		data := make(map[string]interface{})

		db.Fetch(chartDataRows, data)

		chartData = append(chartData, data)

	}

	data["chartData"] = chartData
	return &data, nil
}

func GetChangeFailureRate(params map[string]interface{}) (*map[string]interface{}, errors.Error) {
	medianTimeToRestoreServiceQueryString, formatError := formatQueryString("medianTimeToRestoreServiceQueryString", CHANGE_FAILURE_RATE, params)
	if formatError != nil {
		return nil, errors.Default.Wrap(formatError, "Something went wrong!")
	}
	rows, err := db.RawCursor(medianTimeToRestoreServiceQueryString)

	if err != nil {
		return nil, errors.Default.Wrap(err, "Something went wrong!")
	}

	data := make(map[string]interface{})
	for rows.Next() {
		db.Fetch(rows, data)
	}

	medianTimeToRestoreServiceChartQueryString, formatError := formatQueryString("medianTimeToRestoreServiceChartQueryString", CHANGE_FAILURE_RATE_CHART_QUERY, params)
	if formatError != nil {
		return nil, errors.Default.Wrap(formatError, "Something went wrong!")
	}
	chartDataRows, err := db.RawCursor(medianTimeToRestoreServiceChartQueryString)

	if err != nil {
		return nil, errors.Default.Wrap(err, "Something went wrong!")
	}

	var chartData []interface{}

	for chartDataRows.Next() {
		data := make(map[string]interface{})

		db.Fetch(chartDataRows, data)

		chartData = append(chartData, data)

	}

	data["chartData"] = chartData

	return &data, nil
}
