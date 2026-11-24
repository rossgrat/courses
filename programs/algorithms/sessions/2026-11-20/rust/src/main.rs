fn two_sum(nums: Vec<i32>, target: i32) -> Vec<i32> {
    let mut my_nums = nums.clone();
    my_nums.sort();

    let mut left = 0;
    let mut right = my_nums.len() - 1;

    while left < right {
        let sum = my_nums[left] + my_nums[right];
        if sum == target {
            break;
        } else if sum < target {
            left = left + 1;
        } else {
            right = right - 1;
        }
    }

    let mut ret = Vec::new();

    for i in 0..my_nums.len() {
        if nums[i] == my_nums[left] || nums[i] == my_nums[right] {
            ret.push(i as i32);
            if ret.len() == 2 {
                break;
            }
        }
    }

    return ret;
}

#[test]
fn test_two_sum() {
    assert_eq!(two_sum(vec![5, 10, 22, 13, 4, 6, 2, 1], 10), vec![4, 5]);
    assert_eq!(two_sum(vec![50, 49, 1, 2], 52), vec![0, 3]);
}
