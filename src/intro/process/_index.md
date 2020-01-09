---
menuTitle: Process
title: Filecoin Spec Process (v1)
entries:
- about
- fip
- contributing
- related_resources
---

# 🚀 Pre-launch mode

Until we launch, we are making lots of changes to the spec to finish documenting the current version of the protocol. Changes will be made to the spec by a simple PR process, with approvals by key stakeholders. Some refinements are still to happen and testnet is expected to bring a few significant fixes/improvements. Most changes now are changing _the document_, **NOT** changing _the protocol_, at least not in a major way.

Until we launch, if something is missing, PR it in. If something is wrong, PR a fix. If something needs to be elaborated, PR in updates. What is in the top level of this repo, in master, is the spec, is the Filecoin Protocol. Nothing else matters (ie. no other documents, issues contain "the protocol").

# New Proposals -> Drafts -> Spec

{{% notice warning %}}
⚠️ **WARNING:** Filecoin is in pre-launch mode, and we are finishing protocol spec and implementations of the _current_ construction/version of the protocol only. We are highly unlikely to merge anything new into the Filecoin Protocol until after mainnet. Feel free to explore ideas anyway and prepare improvements for the future.
{{% /notice %}}

For anything that is not part of the currently speced systems (like 'repair', for example) the process we will use is:

- **(1) First, discuss the problem(s) and solution(s) in an issue**
  - Or several issues, if the space is large and multithreaded enough.
  - Work out all the details required to make this proposal work.
- **(2) Write a draft with all the details.**
  - When you feel like a solution is near, write up a draft document that contains all the details, and includes what changes would need to happen to the spec
  - E.g. "Add a System called X with ...", or "Add a library called Y, ...", or "Modify vm/state_tree to include ..."
  - Place this document inside the `src/drafts/` directory.
  - Anybody is welcome to contribute well-reasoned and detailed drafts.
  - (Note: these drafts will give way to FIPs in the future)
- **(3) Seek approval to merge this into the specification.**
  - To seek approval, open an issue and discuss it.
  - If the draft approved by the owners of the filecoin-spec, then the changes to the spec will need to be made in a PR.
  - Once changes make it into the spec, remove the draft.

It is acceptable for a PR for a draft to stay open for quite a while, as thought and discussion on the topic happens. At some point, if the reviewers and the author feel that the current state of the draft is stable enough (though not 'done') then it should be merged into the repo. Further changes to the draft are additional PRs, which may generate more discussion. Comments on these drafts are welcome from anyone, but if you wish to be involved in the actual research process, you will need to devote very considerable time and energy to the process.

# On merging

For anything in the drafts or notes folder, merge yourself after a review from a relevant person. For anything in the top level (canonical spec), @zixuanzh, @anorth, @whyrusleeping or @jbenet will merge after proper review.

# Issues

Issues in the specs repo will be high signal. They will either be proposals, or issues directly relating to problems in the spec. More speculative research questions and discussion will happen in the research repo.
