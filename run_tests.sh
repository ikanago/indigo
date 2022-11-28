set -u

# \033[ is the escape for macOS
color_ok="\033[32;1m"
color_failed="\033[31;1m"
color_off="\033[m"

tmp_dir=/tmp/tmpfs/out
asm_dir=$tmp_dir/asm
bin_dir=$tmp_dir/bin
mkdir -p $tmp_dir $asm_dir $bin_dir

function run_unit_test {
    test_name=$1
    src_file=tests/$test_name/main.go
    asm_file=$asm_dir/$test_name.s
    bin_file=$bin_dir/$test_name

    ./output/indigo tests/$test_name/main.go > $asm_file
    clang -o $bin_file $asm_file

    expected_file=tests/$test_name/expected.txt
    actual_file=$tmp_dir/actual.txt
    $bin_file
    echo $? > $actual_file

    if cmp -s $actual_file $expected_file; then
        echo "${color_ok}[ok]${color_off}     ${test_name}"
    else
        echo "${color_failed}[failed]${color_off} ${test_name}, got" $(cat $actual_file)
    fi
}

function run_all_test {
    for test in ./tests/*; do
        test_name=$(basename $test)
        run_unit_test $test_name
    done
}

run_all_test
