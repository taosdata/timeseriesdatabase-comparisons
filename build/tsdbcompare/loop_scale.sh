for s in {100,200,400,600,800,1000};do
    for i in {5000,2500,1000,1};do
		./write.sh -b $i -w 100 -g 1 -s $s
        for j in {50,16,1};do
            ./write.sh -b $i -w $j -g 0 -s $s
        done
    done
done