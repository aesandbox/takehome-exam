package main

import (
    "context"
    "fmt"
    "github.com/google/go-github/github"
    "golang.org/x/oauth2"
    "log"
    "os"
    "strconv"
    "strings"
)

var ctx = context.Background()
var client *github.Client

// a PR will be mergeable if there is AT LEAST 1 review that has been approved,
// all status checks pass, and it has not already been merged.
func isPRMergeable(reviews []*github.PullRequestReview, pr *github.PullRequest) bool {
    prApproved := false
    for _, review := range reviews {
        if review.GetState() == "APPROVED" {
            prApproved = true
            break
        }
    }
    if !prApproved {
        log.Println("PR is not approved")
        return false
    }

    if pr.GetMergeableState() != "clean" {
        log.Println("PR is not in a mergeable state. Check status check or if already merged")
        return false
    }

    // Prompt user to merge or not
    fmt.Print("Type y or yes to merge, anything else to abort: ")
    var shouldMerge string
    fmt.Scan(&shouldMerge)
    shouldMerge = strings.ToLower(shouldMerge)

    return shouldMerge == "y" || shouldMerge == "yes"
}

func main() {
    if len(os.Args) != 2 {
        log.Fatalf("Usage: %s [Pull Request URL]", os.Args[0])
    }

    urlParts := strings.Split(os.Args[1], "/")
    if len(urlParts) != 7 {
        log.Fatal("You must supply a valid PR URL")
    }

    owner := urlParts[3]
    repo := urlParts[4]

    prNumber, err := strconv.Atoi(urlParts[6])
    if err != nil {
        log.Fatal("The PR URL must end with the PR number")
    }

    token := os.Getenv("GH_ACCESS_TOKEN")
    if token == "" {
        log.Fatal("GH_ACCESS_TOKEN access token must be set")
    }

    tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
    tokenContext := oauth2.NewClient(ctx, tokenSource)
    client = github.NewClient(tokenContext)

    pr, _, err := client.PullRequests.Get(ctx, owner, repo, prNumber)
    if err != nil {
        log.Fatal("Request to get PR info failed. Verify PR URL")
    }

    reviews, _, err := client.PullRequests.ListReviews(ctx, owner, repo, prNumber, nil)
    if err != nil {
        log.Fatal("Request to get reviews failed. Verify PR URL")
    }

    if !isPRMergeable(reviews, pr) {
        log.Fatal("Aborting merge")
    }

    _, _, err = client.PullRequests.Merge(ctx, owner, repo, prNumber, "Merging via script", nil)
    if err != nil {
        log.Fatal("Merge request failed")
    }

    _, err = client.Git.DeleteRef(ctx, owner, repo, "heads/"+pr.GetHead().GetRef())
    if err != nil {
        log.Fatal("Failed to delete PR branch")
    }

    fmt.Println("PR was merged and branch was deleted successfully!")
}
