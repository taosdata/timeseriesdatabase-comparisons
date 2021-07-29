echo "run query test"
# for s in {100,1000};do
# 	./read.sh -w 100 -g 1 -s $s
#    for j in {50,16,1};do
#        ./read.sh -w $j -g 0 -s $s
#    done
# done
for s in {100,1000};do
	./read_2.sh -w 100 -g 1 -s $s -n 'fast'
    for j in {50,16,1};do
        ./read_2.sh -w $j -g 0 -s $s -n 'fast'
    done
done

