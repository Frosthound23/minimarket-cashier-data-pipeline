with date_range as (

    select distinct
        toDate(transaction_date) as date_day
    from {{ ref('stg_transactions') }}

),

date_enriched as (

    select
        date_day,
        toDayOfWeek(date_day) as day_of_week
    from date_range

)

select
    toYYYYMMDD(date_day) as date_key,
    date_day,
    toYear(date_day) as year,
    toQuarter(date_day) as quarter,
    toMonth(date_day) as month,
    formatDateTime(date_day, '%M') as month_name,
    toDayOfMonth(date_day) as day_of_month,
    day_of_week,

    case day_of_week
        when 1 then 'Monday'
        when 2 then 'Tuesday'
        when 3 then 'Wednesday'
        when 4 then 'Thursday'
        when 5 then 'Friday'
        when 6 then 'Saturday'
        when 7 then 'Sunday'
    end as day_name

from date_enriched