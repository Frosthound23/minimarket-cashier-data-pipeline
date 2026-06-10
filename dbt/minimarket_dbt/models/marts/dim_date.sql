{{ config(materialized='table') }}

with date_spine as (
    select distinct
        transaction_date_key as date_day
    from {{ ref('stg_transactions') }}
)

select
    date_day,
    toYear(date_day) as year,
    toMonth(date_day) as month,
    toDayOfMonth(date_day) as day,
    toQuarter(date_day) as quarter,
    toDayOfWeek(date_day) as day_of_week,

    case toDayOfWeek(date_day)
        when 1 then 'Monday'
        when 2 then 'Tuesday'
        when 3 then 'Wednesday'
        when 4 then 'Thursday'
        when 5 then 'Friday'
        when 6 then 'Saturday'
        when 7 then 'Sunday'
    end as day_name,

    case toMonth(date_day)
        when 1 then 'January'
        when 2 then 'February'
        when 3 then 'March'
        when 4 then 'April'
        when 5 then 'May'
        when 6 then 'June'
        when 7 then 'July'
        when 8 then 'August'
        when 9 then 'September'
        when 10 then 'October'
        when 11 then 'November'
        when 12 then 'December'
    end as month_name,

    case
        when toDayOfWeek(date_day) in (6, 7) then true
        else false
    end as is_weekend
from date_spine