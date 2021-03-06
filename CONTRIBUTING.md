# Contributing to Textile

:+1::tada: First off, thanks for taking the time to contribute! :tada::+1:

The following is a set of guidelines for contributing to Textile and its libraries, which are hosted in the [textileio](https://github.com/textileio) organization on GitHub. These are mostly guidelines, not rules. Use your best judgment, and feel free to propose changes to this document in a pull request.

#### Table Of Contents

- [Code of Conduct](#code-of-conduct)
- [I don't want to read this whole thing, I just have a question!!!](#i-dont-want-to-read-this-whole-thing-i-just-have-a-question)
- [How Can I Contribute?](#how-can-i-contribute)
  - [Reporting Bugs](#reporting-bugs)
  - [Suggesting Enhancements](#suggesting-enhancements)
  - [Your First Code Contribution](#your-first-code-contribution)
  - [Pull Requests](#pull-requests)
- [Managing Issues](#managing-issues)
  - [Discussions, Bugs, and Enhancements](#discussions-bugs-and-enhancements)
  - [Ambiguous Issues](#ambiguous-issues)
  - [Abandoned Issues](#abandoned-issues)
- [Styleguides](#styleguides)
  - [Git Commit Messages](#git-commit-messages)
  - [Documentation Styleguide](#documentation-styleguide)
- [Additional Notes](#additional-notes)
  - [Issue and Pull Request Labels](#issue-and-pull-request-labels)

## Code of Conduct

This project and everyone participating in it is governed by the [Contributor Covenant Code of Conduct](./CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to [contact@textile.io](mailto:contact@textile.io).

## I don't want to read this whole thing I just have a question!!!

> **Note:** [Please don't file an issue to ask a question.](#managing-issues) You'll get faster results by using the resources below.

We have an official community message board where the community chimes in with helpful advice if you have questions.

* [The official community message board](https://community.textile.io)

If chat is more your speed, you can join the Community Slack channel:

* [Join the Community Slack channel](https://slack.textile.io/)
    * Even though Slack is a chat service, sometimes it takes several hours for community members to respond &mdash; please be patient!
    * Use the `#general` channel for general questions or discussion related to Textile
    * Use the `#support` channel for specific questions about developing with Textile tools
    * Use the `#introductions` and `#events`channels for those things
    * Follow the `#releases` channel for automated release information for all Textile GitHub repos
    * There are many other channels available, check the channel list

## How Can I Contribute?

### Reporting Bugs

This section guides you through submitting a bug report for Textile. Following these guidelines helps maintainers and the community understand your report :pencil:, reproduce the behavior :computer: :computer:, and find related reports :mag_right:.

Before creating bug reports, please check [this list](#before-submitting-a-bug-report) as you might find out that you don't need to create one. When you are creating a bug report, please [include as many details as possible](#how-do-i-submit-a-good-bug-report). Fill out [the required template](./.github/ISSUE_TEMPLATE/bug_report.md), the information it asks for helps us resolve issues faster.

> **Note:** If you find a **Closed** issue that seems like it is the same thing that you're experiencing, open a new issue and include a link to the original issue in the body of your new one.

#### Before Submitting A Bug Report

* **Check the [community forum](https://community.textile.io)** for some common questions and problems.
* **Determine [which repository the problem should be reported in](https://github.com/textileio)**.
* **Perform a [cursory search](https://github.com/search?q=+is%3Aissue+user%3Atextile)** to see if the problem has already been reported. If it has **and the issue is still open**, add a comment to the existing issue instead of opening a new one.

#### How Do I Submit A (Good) Bug Report?

Bugs are tracked as [GitHub issues](https://guides.github.com/features/issues/). After you've determined [which repository](https://github.com/textileio) your bug is related to, create an issue on that repository and provide the following information by filling in [the template](./.github/ISSUE_TEMPLATE/bug_report.md).

Explain the problem and include additional details to help maintainers reproduce the problem:

* **Use a clear and descriptive title** for the issue to identify the problem.
* **Describe the exact steps which reproduce the problem** in as many details as possible. For example, start by explaining how you started the application, e.g. which command exactly you used in the terminal, or how you started Textile otherwise. When listing steps, **don't just say what you did, but explain how you did it**. For example, if you used a command line tool, make sure you provide the sequence of steps you performed to create the bug in the first place.
* **Provide specific examples to demonstrate the steps**. Include links to files or GitHub projects, or copy/pasteable snippets, which you use in those examples. If you're providing snippets in the issue, use [Markdown code blocks](https://help.github.com/articles/markdown-basics/#multiple-lines).
* **Describe the behavior you observed after following the steps** and point out what exactly is the problem with that behavior.
* **Explain which behavior you expected to see instead and why.**
* **Include screenshots and animated GIFs** which show you following the described steps and clearly demonstrate the problem. You can use [this tool](https://www.cockos.com/licecap/) to record GIFs on macOS and Windows, and [this tool](https://github.com/colinkeenan/silentcast) or [this tool](https://github.com/GNOME/byzanz) on Linux.
* **If you're reporting that something crashed**, include a crash report with a stack trace from the operating system. On macOS, the crash report will be available in `Console.app` under "Diagnostic and usage information" > "User diagnostic reports". Include the crash report in the issue in a [code block](https://help.github.com/articles/markdown-basics/#multiple-lines), a [file attachment](https://help.github.com/articles/file-attachments-on-issues-and-pull-requests/), or put it in a [gist](https://gist.github.com/) and provide link to that gist.
* **If the problem wasn't triggered by a specific action**, describe what you were doing before the problem happened and share more information using the guidelines below.

Provide more context by answering these questions:

* **Did the problem start happening recently** (e.g. after updating to a new version) or was this always a problem?
* If the problem started happening recently, **can you reproduce the problem in an older version?** What's the most recent version in which the problem doesn't happen? You can download older versions from [the releases page](../../releases).
* **Can you reliably reproduce the issue?** If not, provide details about how often the problem happens and under which conditions it normally happens.

Include details about your configuration and environment:

* **Which version of the software are you using?**
* **What's the name and version of the OS you're using**?
* **Are you running in a virtual machine?** If so, which VM software are you using and which operating systems and versions are used for the host and the guest?

### Suggesting Enhancements

This section guides you through submitting an enhancement suggestion for Textile, including completely new features and minor improvements to existing functionality. Following these guidelines helps maintainers and the community understand your suggestion :pencil: and find related suggestions :mag_right:.

Before creating enhancement suggestions, please check [this list](#before-submitting-an-enhancement-suggestion) as you might find out that you don't need to create one. When you are creating an enhancement suggestion, please [include as many details as possible](#how-do-i-submit-a-good-enhancement-suggestion). Fill in [the template](./.github/blob/master/.github/ISSUE_TEMPLATE/feature_request.md), including the steps that you imagine you would take if the feature you're requesting existed.

#### Before Submitting An Enhancement Suggestion

* **Check the [community forum](https://community.textile.io)** for some common questions and problems.
* **Determine [which repository the enhancement should be reported in](https://github.com/textileio)**.
* **Perform a [cursory search](https://github.com/search?q=+is%3Aissue+user%3Atextile)** to see if the enhancement has already been suggested. If it has **and the issue is still open**, add a comment to the existing issue instead of opening a new one.

#### How Do I Submit A (Good) Enhancement Suggestion?

Enhancement suggestions are tracked as [GitHub issues](https://guides.github.com/features/issues/). After you've determined [which repository](https://github.com/textileio) your enhancement suggestion is related to, create an issue on that repository and provide the following information:

* **Use a clear and descriptive title** for the issue to identify the suggestion.
* **Provide a step-by-step description of the suggested enhancement** in as many details as possible.
* **Provide specific examples to demonstrate the steps**. Include copy/pasteable snippets which you use in those examples, as [Markdown code blocks](https://help.github.com/articles/markdown-basics/#multiple-lines).
* **Describe the current behavior** and **explain which behavior you expected to see instead** and why.
* **Include screenshots and animated GIFs** which help you demonstrate the steps or point out what the suggestion is related to. You can use [this tool](https://www.cockos.com/licecap/) to record GIFs on macOS and Windows, and [this tool](https://github.com/colinkeenan/silentcast) or [this tool](https://github.com/GNOME/byzanz) on Linux.
* **Explain why this enhancement would be useful** to most Textile users and isn't something that can or should be implemented elsewhere.
* **List some other tools or applications where this enhancement exists.**
* **Specify which version you're using.**
* **Specify the name and version of the OS you're using.**

### Your First Code Contribution

Unsure where to begin contributing to Textile? You can start by looking through these `beginner` and `help-wanted` issues:

* [Beginner issues][beginner] - issues which should only require a few lines of code, and a test or two.
* [Help wanted issues][help-wanted] - issues which should be a bit more involved than `beginner` issues.

Both issue lists are sorted by total number of comments. While not perfect, number of comments is a reasonable proxy for impact a given change will have.

If you want to read about using Textile or developing on top of Textile, the [Textile Docs](https://docs.textile.io) are a great place to start. You can find the source for the docs in [textileio/docs](https://github.com/textileio/docs).

### Pull Requests

The process described here has several goals:

- Maintain Textile's quality
- Fix problems that are important to users
- Engage the community in working toward the best possible software
- Enable a sustainable system for Textile's maintainers to review contributions

Please follow these steps to have your contribution considered by the maintainers:

1. Follow all instructions in [the template](PULL_REQUEST_TEMPLATE.md)
2. Follow the [styleguides](#styleguides)
3. After you submit your pull request, verify that all [status checks](https://help.github.com/articles/about-status-checks/) are passing <details><summary>What if the status checks are failing?</summary>If a status check is failing, and you believe that the failure is unrelated to your change, please leave a comment on the pull request explaining why you believe the failure is unrelated. A maintainer will re-run the status check for you. If we conclude that the failure was a false positive, then we will open an issue to track that problem with our status check suite.</details>

While the prerequisites above must be satisfied prior to having your pull request reviewed, the reviewer(s) may ask you to complete additional design work, tests, or other changes before your pull request can be ultimately accepted.

## Managing Issues

### Discussions, Bugs, and Enhancements

A lot of items that come in on Issues are not bugs or feature requests, they’re questions or discussions. A lot of these questions can be answered by anyone in the community; they’re not exclusive to a Textile maintainer. For feedback like this we set up our [official Textile Community Forums](https://community.textile.io). It’s a great place where people can get all kinds of help, whether it is about using Textile or just about any IPFS-related topics. For those that like more real-time discussions, there is also the [Community Slack channel](https://slack.textile.io).

The general guide is that GitHub Issues is best for things that have a definition of "done": they can be fixed, added, resolved, have a stake driven through its heart or otherwise laid to rest. For things where there isn't a clear goal or end state, where you want to debate theory, or ask questions, the Focums and Slack are the way to go.

We encouraging people to use the best channel for their required feedback. If you ask a question on Issues and get a simple answer and reminder to check the Forums, it’s because we want you to get the help you need as quickly as possible.

### Ambiguous Issues

The whole point of Issues is that they are things that need to be fixed, implemented or completed in some fashion. But there are classes of things that get reported that are undefined or indistinct, there is no way to complete them or simply are chores that will never be done.

When a report comes in that looks like this, we will ask the original author of the Issue to clarify what "done" would look like to them. We can and will help with this process if you're unsure. But if the goal remains undefined or, in the maintainers' estimation, is unachievable, we will close the Issue.

### Abandoned Issues

Issues are a way for the users and the maintainers of a project to communicate and cooperate towards goals. If the maintainers can't get the information they need to resolve something, an Issue can just sit, get lost in the shuffle and never move forward. We also understand that people have busy lives and sometimes they simply can't get back to us or have forgotten the context of an Issue.

In order to focus the maintainers' attention on making progress, we will mark Issues with the `more-information-needed` label when the maintainers have a question. If we don’t receive a response from the original author within a week, we'll give a gentle reminder. If we still haven't received a response within 30 days, we will close the Issue.

## Styleguides

### Git Commit Messages

* Use the present tense ("Add feature" not "Added feature")
* Use the imperative mood ("Change type to..." not "Changes type to...")
* Limit the first line to 72 characters or less
* Reference issues and pull requests liberally after the first line
* When only changing documentation, include `[ci skip]` in the commit title
* Consider starting the commit message with an applicable emoji:
    * :art: `:art:` when improving the format/structure of the code
    * :racehorse: `:racehorse:` when improving performance
    * :non-potable_water: `:non-potable_water:` when plugging memory leaks
    * :memo: `:memo:` when writing docs
    * :penguin: `:penguin:` when fixing something on Linux
    * :apple: `:apple:` when fixing something on macOS
    * :checkered_flag: `:checkered_flag:` when fixing something on Windows
    * :bug: `:bug:` when fixing a bug
    * :fire: `:fire:` when removing code or files
    * :green_heart: `:green_heart:` when fixing the CI build
    * :white_check_mark: `:white_check_mark:` when adding tests
    * :lock: `:lock:` when dealing with security
    * :arrow_up: `:arrow_up:` when upgrading dependencies
    * :arrow_down: `:arrow_down:` when downgrading dependencies
    * :shirt: `:shirt:` when removing linter warnings

### GoDoc

It is a good idea to use `godoc` when developing documentation for your Go code. Godoc is the way users will consume a package, so it is ok to tweak things to make them more accessible. Run `godoc -http=<hostport>` to start a `godoc` server when developing locally.

Godoc documentation is converted to HTML or to text using the go/doc package, and is based on [this set of rules](https://golang.org/pkg/go/doc/#ToHTML). Godoc also shows example code that is runnable by the testing package, to give you an idea of [the conventions](https://golang.org/pkg/testing/#hdr-Examples). See ["Godoc: documenting Go code"](https://golang.org/doc/articles/godoc_documenting_go_code.html) for how to write good comments for godoc.

### Go Code

[Effective Go](https://golang.org/doc/effective_go.html) is a very good style guide, and all around good resource for writing effective Go code. The Go maintainers also collect [common comments made during reviews of Go code](https://github.com/golang/go/wiki/CodeReviewComments), so that a single detailed explanation can be referred to by shorthands.

* Avoid global variables, even in packages. By doing so you will introduce side effects if the package is included multiple times.
* Use `go fmt` before committing ([Gofmt](https://golang.org/cmd/gofmt/) is a tool that automatically formats Go source code).
* Use `go lint` to lint source code. Golint differs from `gofmt` in that `gofmt` reformats Go source code, whereas `golint` prints out style mistakes.
* Dependencies should be kept to the minimum. The introduction of a new dependency should be argued in any pull request.
* Since Go 1.11, a standard dependency system is available behind the name [Go Modules](https://github.com/golang/go/wiki/Modules). It provides a way to define and lock dependencies for reproducible builds. It should be used whenever possible.
* We should not use any specific library or framework for testing, as the [standard library](https://golang.org/pkg/) provides already everything to get started.




## Additional Notes

### Issue and Pull Request Labels

This section lists the labels we use to help us track and manage issues and pull requests. Most labels are used across all Textile repositories, but some are specific to this one.

[GitHub search](https://help.github.com/articles/searching-issues/) makes it easy to use labels for finding groups of issues or pull requests you're interested in. For example, you might be interested in [open issues across all Textile-owned packages which are labeled as bugs, but still need to be reliably reproduced](https://github.com/search?utf8=%E2%9C%93&q=is%3Aopen+is%3Aissue+user%3Atextileio+label%3Abug+label%3Aneeds-reproduction) or perhaps [open pull requests which haven't been reviewed yet](https://github.com/search?q=is%3Aopen+is%3Apr+user%3Atextileio+comments%3A0). We  encourage you to read about [other search filters](https://help.github.com/articles/searching-issues/) which will help you write more focused queries.

Please open an issue in this repo if you have suggestions for new labels, and if you notice some labels are missing on some repositories, then please open an issue on that repository.

## Thanks

These guidelines are based on the fantastic [contribution guidelines](https://github.com/atom/atom/blob/master/CONTRIBUTING.md) for the Atom editor. Thanks Atom community :+1::tada: