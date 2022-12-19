#!/bin/bash

mc config host add myalias http://object-storage:9000 promotions-user promotions123;
mc mb myalias/promotions;
mc admin user add myalias app test12345;
mc admin policy add myalias app-user-policy ./policy/app-user-policy.json;
mc admin policy set myalias app-user-policy user=app;