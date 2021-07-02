# Performance comparisons between InfluxDB and TDengine
This project is a fork of [InfluxDB comparisions project](https://github.com/influxdata/influxdb-comparisons). The testing methodology and test procedure keeps the same as origin project and we just extend the data loading/quering module to support TDengine format and add serveral query test cases. Detailed testing methodology and procedure please refer to the origin project.

Briefly, this comparision test generates devops data and writes into different format, and loads into the database accordingly, then perform the same queries, finally counts the time consumed.
