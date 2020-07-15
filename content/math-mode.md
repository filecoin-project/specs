---
title: Math Mode
bookhidden: true
math-mode: true
---


{{<plain hidden>}}
$$
\gdef\createporepbatch{\textsf{create_porep_batch}}
\gdef\GrothProof{\textsf{Groth16Proof}}
\gdef\Groth{\textsf{Groth16}}
\gdef\GrothEvaluationKey{\textsf{Groth16EvaluationKey}}
\gdef\GrothVerificationKey{\textsf{Groth16VerificationKey}}
\gdef\creategrothproof{\textsf{create_groth16_proof}}
\gdef\ParentLabels{\textsf{ParentLabels}}
\gdef\or#1#2{\langle #1 | #2 \rangle}
\gdef\porepreplicas{\textsf{porep_replicas}}
\gdef\postreplicas{\textsf{post_replicas}}
\gdef\winningpartitions{\textsf{winning_partitions}}
\gdef\windowpartitions{\textsf{window_partitions}}
\gdef\sector{\textsf{sector}}
\gdef\lebitstolebytes{\textsf{le_bits_to_le_bytes}}
\gdef\lebinrep#1{{\llcorner #1 \lrcorner_{\lower{2pt}{2, \textsf{le}}}}}
\gdef\bebinrep#1{{\llcorner #1 \lrcorner_{\lower{2pt}{2, \textsf{be}}}}}
\gdef\lebytesbinrep#1{{\llcorner #1 \lrcorner_{\lower{2pt}{2, \textsf{le-bytes}}}}}
\gdef\fesitelrounds{\textsf{fesitel_rounds}}
\gdef\int{\textsf{int}}
\gdef\lebytes{\textsf{le-bytes}}
\gdef\lebytestolebits{\textsf{le_bytes_to_le_bits}}
\gdef\letooctet{\textsf{le_to_octet}}
\gdef\byte{\textsf{byte}}
\gdef\postpartitions{\textsf{post_partitions}}
\gdef\PostReplica{\textsf{PostReplica}}
\gdef\PostReplicas{\textsf{PostReplicas}}
\gdef\PostPartitionProof{\textsf{PostPartitionProof}}
\gdef\PostReplicaProof{\textsf{PostReplicaProof}}
\gdef\TreeRProofs{\textsf{TreeRProofs}}
\gdef\pad{\textsf{pad}}
\gdef\octettole{\textsf{octet_to_le}}
\gdef\packed{\textsf{packed}}
\gdef\val{\textsf{val}}
\gdef\bits{\textsf{bits}}
\gdef\partitions{\textsf{partitions}}
\gdef\Batch{\textsf{Batch}}
\gdef\batch{\textsf{batch}}
\gdef\postbatch{\textsf{post_batch}}
\gdef\postchallenges{\textsf{post_challenges}}
\gdef\Nonce{\textsf{Nonce}}
\gdef\createvanillaporepproof{\textsf{create_vanilla_porep_proof}}
\gdef\PorepVersion{\textsf{PorepVersion}}
\gdef\bedecode{\textsf{be_decode}}
\gdef\OR{\mathbin{|}}
\gdef\indexbits{\textsf{index_bits}}
\gdef\nor{\textsf{nor}}
\gdef\and{\textsf{and}}
\gdef\norgadget{\textsf{nor_gadget}}
\gdef\andgadget{\textsf{and_gadget}}
\gdef\el{\textsf{el}}
\gdef\arr{\textsf{arr}}
\gdef\pickgadget{\textsf{pick_gadget}}
\gdef\pick{\textsf{pick}}
\gdef\int{\textsf{int}}
\gdef\x{\textsf{x}}
\gdef\y{\textsf{y}}
\gdef\aap{{\langle \auxb | \pubb \rangle}}
\gdef\aapc{{\langle \auxb | \pubb | \constb \rangle}}
\gdef\TreeRProofs{\textsf{TreeRProofs}}
\gdef\parentlabelsbits{\textsf{parent_labels_bits}}
\gdef\label{\textsf{label}}
\gdef\layerbits{\textsf{layer_bits}}
\gdef\labelbits{\textsf{label_bits}}
\gdef\digestbits{\textsf{digest_bits}}
\gdef\node{\textsf{node}}
\gdef\layerindex{\textsf{layer_index}}
\gdef\be{\textsf{be}}
\gdef\octet{\textsf{octet}}
\gdef\reverse{\textsf{reverse}}
\gdef\LSBit{\textsf{LSBit}}
\gdef\MSBit{\textsf{MSBit}}
\gdef\LSByte{\textsf{LSByte}}
\gdef\MSByte{\textsf{MSByte}}
\gdef\PorepPartitionProof{\textsf{PorepPartitionProof}}
\gdef\PostPartitionProof{\textsf{PostPartitionProof}}
\gdef\lebinrep#1{{\llcorner #1 \lrcorner_{\lower{2pt}{2, \textsf{le}}}}}
\gdef\bebinrep#1{{\llcorner #1 \lrcorner_{\lower{2pt}{2, \textsf{be}}}}}
\gdef\octetbinrep#1{{\llcorner #1 \lrcorner_{\lower{2pt}{2, \textsf{octet}}}}}
\gdef\fieldelement{\textsf{field_element}}
\gdef\Fqsafe{{\mathbb{F}_{q, \safe}}}
\gdef\elem{\textsf{elem}}
\gdef\challenge{\textsf{challenge}}
\gdef\challengeindex{\textsf{challenge_index}}
\gdef\uniquechallengeindex{\textsf{unique_challenge_index}}
\gdef\replicaindex{\textsf{replica_index}}
\gdef\uniquereplicaindex{\textsf{unique_replica_index}}
\gdef\nreplicas{\textsf{n_replicas}}
\gdef\unique{\textsf{unique}}
\gdef\R{\mathcal{R}}
\gdef\getpostchallenge{\textsf{get_post_challenge}}
\gdef\verifyvanillapostproof{\textsf{verify_vanilla_post_proof}}
\gdef\BinPathElement{\textsf{BinPathElement}}
\gdef\BinTreeDepth{\textsf{BinTreeDepth}}
\gdef\BinTree{\textsf{BinTree}}
\gdef\BinTreeProof{\textsf{BinTreeProof}}
\gdef\bintreeproofisvalid{\textsf{bintree_proof_is_valid}}
\gdef\Bit{{\{0, 1\}}}
\gdef\Byte{\mathbb{B}}
\gdef\calculatebintreechallenge{\textsf{calculate_bintree_challenge}}
\gdef\calculateocttreechallenge{\textsf{calculate_octtree_challenge}}
\gdef\depth{\textsf{depth}}
\gdef\dot{\textsf{.}}
\gdef\for{\textsf{for }}
\gdef\Function{\textbf{Function: }}
\gdef\Fq{{\mathbb{F}_q}}
\gdef\leaf{\textsf{leaf}}
\gdef\line#1#2#3{\scriptsize{\textsf{#1.}#2}\ \normalsize{#3}}
\gdef\missing{\textsf{missing}}
\gdef\NodeIndex{\textsf{NodeIndex}}
\gdef\nodes{\textsf{nodes}}
\gdef\OctPathElement{\textsf{OctPathElement}}
\gdef\OctTree{\textsf{OctTree}}
\gdef\OctTreeDepth{\textsf{OctTreeDepth}}
\gdef\OctTreeProof{\textsf{OctTreeProof}}
\gdef\octtreeproofisvalid{\textsf{octtree_proof_is_valid}}
\gdef\path{\textsf{path}}
\gdef\pathelem{\textsf{path_elem}}
\gdef\return{\textsf{return }}
\gdef\root{\textsf{root}}
\gdef\Safe{{\Byte^{[32]}_\textsf{safe}}}
\gdef\sibling{\textsf{sibling}}
\gdef\siblings{\textsf{siblings}}
\gdef\struct{\textsf{struct }}
\gdef\Teq{\underset{{\small \mathbb{T}}}{=}}
\gdef\Tequiv{\underset{{\small \mathbb{T}}}{\equiv}}
\gdef\thin{{\thinspace}}
\gdef\AND{\mathbin{\&}}
\gdef\MOD{\mathbin{\%}}
\gdef\createproof{{\textsf{create\_proof}}}
\gdef\layer{\textsf{layer}}
\gdef\nodeindex{\textsf{node_index}}
\gdef\childindex{\textsf{child_index}}
\gdef\push{\textsf{push}}
\gdef\index{\textsf{index}}
\gdef\leaves{\textsf{leaves}}
\gdef\len{\textsf{len}}
\gdef\ColumnProof{\textsf{ColumnProof}}
\gdef\concat{\ \|\ }
\gdef\inputs{\textsf{inputs}}
\gdef\Poseidon{\textsf{Poseidon}}
\gdef\bi{\ \ }
\gdef\Bool{{\{\textsf{True}, \textsf{False}\}}}
\gdef\curr{\textsf{curr}}
\gdef\if{\textsf{if }}
\gdef\else{\textsf{else}}
\gdef\proof{\textsf{proof}}
\gdef\Sha#1{\textsf{Sha#1}}
\gdef\ldotdot{{\ldotp\ldotp}}
\gdef\as{\textsf{ as }}
\gdef\bintreerootgadget{\textsf{bintree_root_gadget}}
\gdef\octtreerootgadget{\textsf{octtree_root_gadget}}
\gdef\cs{\textsf{cs}}
\gdef\RCS{\textsf{R1CS}}
\gdef\pathbits{\textsf{path_bits}}
\gdef\missingbit{\textsf{missing_bit}}
\gdef\missingbits{\textsf{missing_bits}}
\gdef\pubb{\textbf{pub}}
\gdef\privb{\textbf{priv}}
\gdef\auxb{\textbf{aux}}
\gdef\constb{\textbf{const}}
\gdef\CircuitVal{\textsf{CircuitVal}}
\gdef\CircuitBit{{\textsf{CircuitVal}_\Bit}}
\gdef\Le{\textsf{le}}
\gdef\privateinput{\textsf{private_input}}
\gdef\publicinput{\textsf{public_input}}
\gdef\deq{\mathbin{\overset{\diamond}{=}}}
\gdef\alloc{\textsf{alloc}}
\gdef\insertgadget#1{\textsf{insert_#1_gadget}}
\gdef\block{\textsf{block}}
\gdef\shagadget#1#2{\textsf{sha#1_#2_gadget}}
\gdef\poseidongadget#1{\textsf{poseidon_#1_gadget}}
\gdef\refeq{\mathbin{\overset{{\small \&}}=}}
\gdef\ptreq{\mathbin{\overset{{\small \&}}=}}
\gdef\bit{\textsf{bit}}
\gdef\auxle{{[\textbf{aux}, \textsf{le}]}}
\gdef\SpecificNotation{{\underline{\text{Specific Notation}}}}
\gdef\repeat{\textsf{repeat}}
\gdef\preimage{\textsf{preimage}}
\gdef\digest{\textsf{digest}}
\gdef\digestbytes{\textsf{digest_bytes}}
\gdef\digestint{\textsf{digest_int}}
\gdef\leencode{\textsf{le_encode}}
\gdef\ledecode{\textsf{le_decode}}
\gdef\ReplicaID{\textsf{ReplicaID}}
\gdef\replicaid{\textsf{replica_id}}
\gdef\replicaidbits{\textsf{replica_id_bits}}
\gdef\replicaidblock{\textsf{replica_id_block}}
\gdef\cc{\textsf{::}}
\gdef\new{\textsf{new}}
\gdef\lebitsgadget{\textsf{le_bits_gadget}}
\gdef\CircuitBitOrConst{{\textsf{CircuitValOrConst}_\Bit}}
\gdef\createporepcircuit{\textsf{create_porep_circuit}}
\gdef\CommD{\textsf{CommD}}
\gdef\CommC{\textsf{CommC}}
\gdef\CommR{\textsf{CommR}}
\gdef\CommCR{\textsf{CommCR}}
\gdef\commd{\textsf{comm_d}}
\gdef\commc{\textsf{comm_c}}
\gdef\commr{\textsf{comm_r}}
\gdef\commcr{\textsf{comm_cr}}
\gdef\assert{\textsf{assert}}
\gdef\asserteq{\textsf{assert_eq}}
\gdef\TreeDProof{\textsf{TreeDProof}}
\gdef\TreeRProof{\textsf{TreeRProof}}
\gdef\TreeR{\textsf{TreeR}}
\gdef\ParentColumnProofs{\textsf{ParentColumnProofs}}
\gdef\challengebits{\textsf{challenge_bits}}
\gdef\packedchallenge{\textsf{packed_challenge}}
\gdef\PartitionProof{\textsf{PartitionProof}}
\gdef\u#1{\textsf{u#1}}
\gdef\packbitsasinputgadget{\textsf{pack_bits_as_input_gadget}}
\gdef\treedleaf{\textsf{tree_d_leaf}}
\gdef\treerleaf{\textsf{tree_r_leaf}}
\gdef\calculatedtreedroot{\textsf{calculated_tree_d_root}}
\gdef\calculatedtreerleaf{\textsf{calculated_tree_r_leaf}}
\gdef\calculatedcommd{\textsf{calculated_comm_d}}
\gdef\calculatedcommc{\textsf{calculated_comm_c}}
\gdef\calculatedcommr{\textsf{calculated_comm_r}}
\gdef\calculatedcommcr{\textsf{calculated_comm_cr}}
\gdef\layers{\textsf{layers}}
\gdef\total{\textsf{total}}
\gdef\column{\textsf{column}}
\gdef\parentcolumns{\textsf{parent_columns}}
\gdef\columns{\textsf{columns}}
\gdef\parentlabel{\textsf{parent_label}}
\gdef\label{\textsf{label}}
\gdef\calculatedtreecleaf{\textsf{calculated_tree_c_leaf}}
\gdef\calculatedcolumn{\textsf{calculated_column}}
\gdef\parentlabels{\textsf{parent_labels}}
\gdef\drg{\textsf{drg}}
\gdef\exp{\textsf{exp}}
\gdef\parentlabelbits{\textsf{parent_label_bits}}
\gdef\parentlabelblock{\textsf{parent_label_block}}
\gdef\Bits{\textsf{ Bits}}
\gdef\safe{\textsf{safe}}
\gdef\calculatedlabel{\textsf{calculated_label}}
\gdef\createlabelgadget{\textsf{create_label_gadget}}
\gdef\encodingkey{\textsf{encoding_key}}
\gdef\encodegadget{\textsf{encode_gadget}}
\gdef\TreeC{\textsf{TreeC}}
\gdef\value{\textsf{value}}
\gdef\encoded{\textsf{encoded}}
\gdef\unencoded{\textsf{unencoded}}
\gdef\key{\textsf{key}}
\gdef\lc{\textsf{lc}}
\gdef\LC{\textsf{LC}}
\gdef\LinearCombination{\textsf{LinearCombination}}
\gdef\one{\textsf{one}}
\gdef\constraint{\textsf{constraint}}
\gdef\proofs{\textsf{proofs}}
\gdef\merkleproofs{\textsf{merkle_proofs}}
\gdef\TreeRProofs{\textsf{TreeRProofs}}
\gdef\challenges{\textsf{challenges}}
\gdef\pub{\textsf{pub}}
\gdef\priv{\textsf{priv}}
\gdef\last{\textsf{last}}
\gdef\TreeRProofs{\textsf{TreeRProofs}}
\gdef\post{\textsf{post}}
\gdef\SectorID{\textsf{SectorID}}
\gdef\winning{\textsf{winning}}
\gdef\window{\textsf{window}}
\gdef\Replicas{\textsf{Replicas}}
\gdef\P{\mathcal{P}}
\gdef\ww{{\textsf{winning}|\textsf{window}}}
\gdef\replicasperk{{\textsf{replicas}/k}}
\gdef\replicas{\textsf{replicas}}
\gdef\Replica{\textsf{Replica}}
\gdef\createvanillapostproof{\textsf{create_vanilla_post_proof}}
\gdef\createpostcircuit{\textsf{create_post_circuit}}
\gdef\ReplicaProof{\textsf{ReplicaProof}}
\gdef\aww{{\langle \ww \rangle}}
\gdef\partitionproof{\textsf{partition_proof}}
\gdef\replicas{\textsf{replicas}}
\gdef\getdrgparents{\textsf{get_drg_parents}}
\gdef\getexpparents{\textsf{get_exp_parents}}
\gdef\DrgSeed{\textsf{DrgSeed}}
\gdef\DrgSeedPrefix{\textsf{DrgSeedPrefix}}
\gdef\FeistelKeysBytes{\textsf{FeistelKeysBytes}}
\gdef\porep{\textsf{porep}}
\gdef\rng{\textsf{rng}}
\gdef\ChaCha#1{\textsf{ChaCha#1}}
\gdef\cc{\textsf{::}}
\gdef\fromseed{\textsf{from_seed}}
\gdef\buckets{\textsf{buckets}}
\gdef\meta{\textsf{meta}}
\gdef\dist{\textsf{dist}}
\gdef\each{\textsf{each}}
\gdef\PorepID{\textsf{PorepID}}
\gdef\porepgraphseed{\textsf{porep_graph_seed}}
\gdef\utf{\textsf{utf8}}
\gdef\DrgStringID{\textsf{DrgStringID}}
\gdef\FeistelStringID{\textsf{FeistelStringID}}
\gdef\graphid{\textsf{graph_id}}
\gdef\createfeistelkeys{\textsf{create_feistel_keys}}
\gdef\FeistelKeys{\textsf{FeistelKeys}}
\gdef\feistelrounds{\textsf{fesitel_rounds}}
\gdef\feistel{\textsf{feistel}}
\gdef\ExpEdgeIndex{\textsf{ExpEdgeIndex}}
\gdef\loop{\textsf{loop}}
\gdef\right{\textsf{right}}
\gdef\left{\textsf{left}}
\gdef\mask{\textsf{mask}}
\gdef\RightMask{\textsf{RightMask}}
\gdef\LeftMask{\textsf{LeftMask}}
\gdef\roundkey{\textsf{round_key}}
\gdef\beencode{\textsf{be_encode}}
\gdef\Blake{\textsf{Blake2b}}
\gdef\input{\textsf{input}}
\gdef\output{\textsf{output}}
\gdef\while{\textsf{while }}
\gdef\digestright{\textsf{digest_right}}
\gdef\xor{\mathbin{\oplus_\text{xor}}}
\gdef\Edges{\textsf{ Edges}}
\gdef\edge{\textsf{edge}}
\gdef\expedge{\textsf{exp_edge}}
\gdef\expedges{\textsf{exp_edges}}
\gdef\createlabel{\textsf{create_label}}
\gdef\Label{\textsf{Label}}
\gdef\Column{\textsf{Column}}
\gdef\Columns{\textsf{Columns}}
\gdef\ParentColumns{\textsf{ParentColumns}}
%\gdef\tern#1?#2:#3{#1\ \text{?}\ #2 \ \text{:}\ #3}
\gdef\repeattolength{\textsf{repeat_to_length}}
\gdef\verifyvanillaporepproof{\textsf{verify_vanilla_porep_proof}}
\gdef\poreppartitions{\textsf{porep_partitions}}
\gdef\challengeindex{\textsf{challenge_index}}
\gdef\porepbatch{\textsf{porep_batch}}
\gdef\winningchallenges{\textsf{winning_challenges}}
\gdef\windowchallenges{\textsf{window_challenges}}
\gdef\PorepPartitionProof{\textsf{PorepPartitionProof}}
\gdef\TreeD{\textsf{TreeD}}
\gdef\TreeCProof{\textsf{TreeCProof}}
\gdef\Labels{\textsf{Labels}}
\gdef\porepchallenges{\textsf{porep_challenges}}
\gdef\postchallenges{\textsf{post_challenges}}
\gdef\PorepChallengeSeed{\textsf{PorepChallengeSeed}}
\gdef\getporepchallenges{\textsf{get_porep_challenges}}
\gdef\getallparents{\textsf{get_all_parents}}
\gdef\PorepChallengeProof{\textsf{PorepChallengeProof}}
\gdef\challengeproof{\textsf{challenge_proof}}
\gdef\PorepChallenges{\textsf{PorepChallenges}}
\gdef\replicate{\textsf{replicate}}
\gdef\createreplicaid{\textsf{create_replica_id}}
\gdef\ProverID{\textsf{ProverID}}
\gdef\replicaid{\textsf{replica_id}}
\gdef\generatelabels{\textsf{generate_labels}}
\gdef\labelwithdrgparents{\textsf{label_with_drg_parents}}
\gdef\labelwithallparents{\textsf{label_with_all_parents}}
\gdef\createtreecfromlabels{\textsf{create_tree_c_from_labels}}
\gdef\ColumnDigest{\textsf{ColumnDigest}}
\gdef\encode{\textsf{encode}}
$$
{{</plain>}}

# Math mode
---

## SDR Spec

### Merkle Proofs

**Implementation:**
* [`storage_proofs::merkle::MerkleTreeWrapper::gen_proof()`]()
* [`merkle_light::merkle::MerkleTree::gen_proof()`](https://github.com/filecoin-project/merkle_light/blob/64a468807c594d306d12d943dd90cc5f88d0d6b0/src/merkle.rs#L918)

**Additional Notation:**
`$\index_l: [\lfloor N_\nodes / 2^l \rfloor] \equiv [\len(\BinTree\dot\layer_l)]$`
The index of a node in a `$\BinTree$` layer `$l$`. The leftmost node in a tree has `$\index_l = 0$`. For each tree layer `$l$` (excluding the root layer) a Merkle proof verifier calculates the label of the node at `$\index_l$` from a single Merkle proof path element `$\BinTreeProof_c\dot\path[l - 1] \thin$`.

### BinTreeProofs

The method `$\BinTreeProof\dot\createproof$` is used to create a Merkle proof for a challenge node `$c$`.

```text
$\overline{\underline{\Function \BinTree\dot\createproof(c: \NodeIndex) \rightarrow \BinTreeProof_c}}$
$\line{1}{\bi}{\leaf: \Safe = \BinTree\dot\leaves[c]}$
$\line{2}{\bi}{\root: \Safe = \BinTree\dot\root}$

$\line{3}{\bi}{\path: \BinPathElement^{[\BinTreeDepth]}= [\ ]}$
$\line{4}{\bi}{\for l \in [\BinTreeDepth]:}$
$\line{5}{\bi}{\quad \index_l: [\len(\BinTree\dot\layer_l)] = c \gg l}$
$\line{6}{\bi}{\quad \missing: \Bit = \index_l \AND 1}$
$\line{7}{\bi}{\quad \sibling: \Safe = \if \missing = 0:}$
$\quad\quad\quad \BinTree\dot\layer_l[\index_l + 1]$
$\quad\quad\thin \else:$
$\quad\quad\quad \BinTree\dot\layer_l[\index_l - 1]$
$\line{8}{\bi}{\quad \path\dot\push(\BinPathElement \thin \{\ \sibling, \thin \missing\ \} \thin )}$

$\line{9}{\bi}{\return \BinTreeProof_c \thin \{\ \leaf, \thin \root, \thin \path\ \}}$
```


**Code Comments:**
* **Line 5:** Calculates the node index in layer `$l$` of the node that the verifier calculated using the previous lath element (or the `$\BinTreeProof_c\dot\leaf$ if $l = 0$`). Note that `$c \gg l \equiv \lfloor c / 2^l \rfloor \thin$`.

### OctTreeProofs

The method $\OctTreeProof\dot\createproof$ is used to create a Merkle proof for a challenge node $c$.

**Additional Notation:**
`$\index_l: [\lfloor N_\nodes / 8^l \rfloor] \equiv [\len(\OctTree\dot\layer_l)]$`
The index of a node in an `$\OctTree$` layer `$l$`. The leftmost node in a tree has `$\index_l = 0$`. For each tree layer `$l$` (excluding the root layer) a Merkle proof verifier calculates the label of the node at `$\index_l$` from a single Merkle proof path element `$\OctTreeProof_c\dot\path[l - 1] \thin$`.

`$\textsf{first\_sibling}_l \thin, \textsf{last\_sibling}: [\lfloor N_\nodes / 8^l \rfloor]$`
The node indexes in tree layer `$l$` of the first and last nodes in this layer's Merkle path element's siblings array `$\OctTreeProof_c\dot\path[l]\dot\siblings \thin$`.

```text
$\overline{\underline{\Function \OctTree\dot\createproof(c: \NodeIndex) \rightarrow \OctTreeProof_c}}$
$\line{1}{\bi}{\leaf: \Fq = \OctTree\dot\leaves[c]}$
$\line{2}{\bi}{\root: \Fq = \OctTree\dot\root}$

$\line{3}{\bi}{\path: \OctPathElement^{[\OctTreeDepth]}= [\ ]}$
$\line{4}{\bi}{\for l \in [\OctTreeDepth]:}$
$\line{5}{\bi}{\quad \index_l: [\len(\OctTree\dot\layer_l)] = c \gg (3 * l)}$
$\line{6}{\bi}{\quad \missing: [8] = \index_l \MOD 8}$

$\line{7}{\bi}{\quad \textsf{first\_sibling}_l = \index_l - \missing}$
$\line{8}{\bi}{\quad \textsf{last\_sibling}_l = \index_l + (7 - \missing)}$
$\line{9}{\bi}{\quad \siblings: \Fq^{[7]} =}$
$\quad\quad\quad \OctTree\dot\layer_l[\textsf{first\_sibling}_l \mathbin{\ldotdot} \index_l]$
$\quad\quad\quad \|\ \OctTree\dot\layer_l[\index_l + 1 \mathbin{\ldotdot} \textsf{last\_sibling}_l + 1]$

$\line{10}{}{\quad \path\dot\push(\OctPathElement \thin \{\ \siblings, \thin \missing\ \} \thin )}$
$\line{11}{}{\return \OctTreeProof_c \thin \{\ \leaf, \thin \root, \thin \path\ \}}$
```