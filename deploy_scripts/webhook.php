<?php

//$cmd = "echo 'Hello, world!' > log.txt";

$cmd = "bash deploy_docker.sh > log_deploy_docker.txt 2>&1";

exec($cmd);

echo "complete";