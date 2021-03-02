for i in {5000,2500,1000,1};do
	for j in {100,50,16,1};do
		./write.sh -b $i -w $j -g 0
	done
done
