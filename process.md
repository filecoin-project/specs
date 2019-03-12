# Spec Process (v0)

## 'Catch Up' Mode

Until we get to 'spec parity' where our current level of understanding of the protocol and the spec are in sync, changes will be made to the spec by a simple PR process. If something is missing, PR it in, if something is wrong, PR a fix, if something needs to be elaborated, PR in updates. What is in the top level of this repo, in master, is the spec.

## Proposals -> Drafts -> Spec

For anything that is not 'catching up' (like 'repair', for example) the process we will use is to first discuss the problem in an issue (or several issues, if the space is large and multithreaded enough). Then when someone feels like a solution is near, they will write it up as a document, and submit a PR to put it into the 'drafts' folder in the repo.

'Drafts' are not canonical spec, and should not be considered for implementation. It is acceptable for a PR for a draft to stay open for quite a while, as thought and discussion on the topic happens. At some point (ideally, in two weeks or less), if the reviewers and the author feel that the current state of the draft is stable enough (though not 'done') then it should be merged into the repo. Further changes to the draft are additional PRs, which may generate more discussion. Comments on these drafts are welcome from anyone, but if you wish to be involved in the actual research process, you will need to devote very considerable time and energy to the process.

Once there is agreement that the draft should be implemented, it should then get moved from the drafts folder, into the top level along with other spec documents. This process should just be a simple renaming, and should not generate any discussion. Along with the moving of that document, any interested parties in the development teams should be explicitly tagged.

### On merging

For anything in the drafts or notes folder, merge yourself after a review from a relevant person. For anything in the top level (canonical spec), @whyrusleeping will merge after proper review.

### Issues

Issues in the specs repo will be high signal, they will either be proposals, or issues directly relating to problems in the spec. More speculative research questions and discussion will happen in the research repo.
