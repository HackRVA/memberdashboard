from ubuntu

RUN mkdir /backup
# create a working directory
WORKDIR /backup

# add backup script
ADD backup.sh .

RUN apt-get update && \
  apt-get install -y \ 
  cron

# Add crontab file in the cron directory
RUN touch /etc/cron.d/db-backup
# # Give execution rights on the cron job
RUN chmod 0644 /etc/cron.d/db-backup
# Setup cronjob to run everyday at 1am
RUN echo '0 1 * * * cd /backup && sh backup.sh' > /etc/cron.d/db-backup
# # Apply cron job
RUN crontab /etc/cron.d/db-backup

CMD cron -f
