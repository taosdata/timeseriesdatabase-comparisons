for s in {100,200,400,600,800,1000};do
    for i in {5000,2500,1000,1};do
		./write_to_server.sh -b $i -w 100 -g 1 -s $s -a "test217" -T "/var/lib/taos/"  -t '2018-01-01T00:00:00Z' -e '2018-01-01T00:10:00Z'
        for j in {50,16,1};do
            ./write_to_server.sh -b $i -w $j -g 0 -s $s -a "test217" -I "/var/lib/influxdb/"  -t '2018-01-01T00:00:00Z' -e '2018-01-01T00:10:00Z'
        done
    done
done