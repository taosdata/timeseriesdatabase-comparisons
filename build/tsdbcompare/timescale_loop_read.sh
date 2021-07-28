for s in {100,1000};do
		./timescale_read.sh -w 1 -g 1 -s $s
        for j in {16,50,100};do
            ./timescale_read.sh  -w $j -g 0 -s $s
        done
done
