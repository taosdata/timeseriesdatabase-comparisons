for s in {100, 1000};do
	./read.sh -w 100 -g 1 -s $s
    for j in {50,16,1};do
        ./write_to_server.sh -w $j -g 0 -s $s
    done
done