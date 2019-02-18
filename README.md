# takehome-exam #

This is a small go program to merge and clean up pull requests, given a pull request URL.

## Usage ##
If you clone this repo as is, and run `go build`:
```sh
./takehome-exam [PullRequestURL]
```

The program will prompt you to merge if the following conditions are met:
* A valid pull request URL is the first, and only argument.
* The most recent commit has an approved review.
* All status checks have passed.

When these conditions are not met, you will be presented with an error saying what was missed. Otherwise, you will see:
`Type y or yes to merge, anything else to abort:`

If the merge and branch deletion are successful, you will see: `PR was merged and branch was deleted successfully!`.
Otherwise, the program will exit with a non-zero exit code and print an error message telling you what is wrong.

## TODO List ##
* Pull input validation out of main() and into its' own function.
* In `func isPRMergeable(...)`:
    * Loop through reviews in reverse, as the most recent reviews are last in the slice.
    * Potentially combine `!prApproved` and `pr.GetMergableState` if statements to reduce number of return statements. It depends on how specific we want error messages to be.
    * Enhance input URL validation.
* Prompt user and request a review if the most recent commit has no reviews.
* Add a flag to bypass prompt asking if you want to merge.
* Add a flag to pass a custom commit message. Currently the message will always be: `Merging via script`
* Add a flag to request a review if there are no approved reviews.


