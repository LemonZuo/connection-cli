#!/bin/sh

# Check log file existence
if [ ! -f /app/app.log ]; then
    touch /app/app.log
fi

# Print Version Information
/app/connection-cli --version

# Print Initialization Message
echo "Connection CLI is running in daemon mode." >> /app/app.log
echo "To run tests, set the MODE environment variable to one of: mysql, postgres, redis, port, http" >> /app/app.log
echo "Example: docker run -e MODE=port -e HOST=example.com -e PORT=80 connection-cli" >> /app/app.log

# Output the log file to the console
tail -f /app/app.log
