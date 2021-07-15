for s in {100,200,400,600,800,1000};do
    for i in {5000,2500,1000};do
		./timescale_write_to_server.sh -b $i -w 100 -g 1 -s $s
        for j in {50,16,8};do
            ./timescale_write_to_server.sh -b $i -w $j -g 0 -s $s
        done
    done
done
