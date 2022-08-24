# Database

The database runs Postgres.

I recommend using pgadmin to interact with the database.  Ideally direct manipulation of the DB should be avoided (there are some edge cases where this is necessary).
The database schema is managed with migrations.  see the main readme of this repo for more info.


## Tables

| table name | description |
| ----- | ----- |
| access_events | We store access events here.  This currently isn't used by anything other than reports (e.g. how many swipes per day). |
| *communication | The communication table is basically an enum of types of messages that we can send out |
| *communication_log | a log of messages that we have sent out |
| member_counts | Everyday, we update how many members we have for each membership level. This allows us to track how our membership has changed each month |
| member_credit | deprecated - can be removed |
| member_resource | stores the relationship between members and what resources they have access to |
| member_tiers | an enum of member levels |
| members | membership information |
| payments | deprecated - can be removed - will be removed on next migrate up |
| resources | resource information - name, address, how to communicate with the resource |
| users | users are tied to members.  The distinction is that users use the dashboard.  We don't support non-members making user accounts |


> *Emails are currently disabled