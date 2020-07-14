$$
\def\createporepbatch{\textsf{create_porep_batch}}
\def\GrothProof{\textsf{Groth16Proof}}
\def\Groth{\textsf{Groth16}}
\def\GrothEvaluationKey{\textsf{Groth16EvaluationKey}}
\def\GrothVerificationKey{\textsf{Groth16VerificationKey}}
\def\creategrothproof{\textsf{create_groth16_proof}}
\def\ParentLabels{\textsf{ParentLabels}}
\def\or#1#2{\langle #1 | #2 \rangle}
\def\porepreplicas{\textsf{porep_replicas}}
\def\postreplicas{\textsf{post_replicas}}
\def\winningpartitions{\textsf{winning_partitions}}
\def\windowpartitions{\textsf{window_partitions}}
\def\sector{\textsf{sector}}
\def\lebitstolebytes{\textsf{le_bits_to_le_bytes}}
\def\lebinrep#1{{\llcorner #1 \lrcorner_{\lower{2pt}{2, \textsf{le}}}}}
\def\bebinrep#1{{\llcorner #1 \lrcorner_{\lower{2pt}{2, \textsf{be}}}}}
\def\lebytesbinrep#1{{\llcorner #1 \lrcorner_{\lower{2pt}{2, \textsf{le-bytes}}}}}
\def\fesitelrounds{\textsf{fesitel_rounds}}
\def\int{\textsf{int}}
\def\lebytes{\textsf{le-bytes}}
\def\lebytestolebits{\textsf{le_bytes_to_le_bits}}
\def\letooctet{\textsf{le_to_octet}}
\def\byte{\textsf{byte}}
\def\postpartitions{\textsf{post_partitions}}
\def\PostReplica{\textsf{PostReplica}}
\def\PostReplicas{\textsf{PostReplicas}}
\def\PostPartitionProof{\textsf{PostPartitionProof}}
\def\PostReplicaProof{\textsf{PostReplicaProof}}
\def\TreeRProofs{\textsf{TreeRProofs}}
\def\pad{\textsf{pad}}
\def\octettole{\textsf{octet_to_le}}
\def\packed{\textsf{packed}}
\def\val{\textsf{val}}
\def\bits{\textsf{bits}}
\def\partitions{\textsf{partitions}}
\def\Batch{\textsf{Batch}}
\def\batch{\textsf{batch}}
\def\postbatch{\textsf{post_batch}}
\def\postchallenges{\textsf{post_challenges}}
\def\Nonce{\textsf{Nonce}}
\def\createvanillaporepproof{\textsf{create_vanilla_porep_proof}}
\def\PorepVersion{\textsf{PorepVersion}}
\def\bedecode{\textsf{be_decode}}
\def\OR{\mathbin{|}}
\def\indexbits{\textsf{index_bits}}
\def\nor{\textsf{nor}}
\def\and{\textsf{and}}
\def\norgadget{\textsf{nor_gadget}}
\def\andgadget{\textsf{and_gadget}}
\def\el{\textsf{el}}
\def\arr{\textsf{arr}}
\def\pickgadget{\textsf{pick_gadget}}
\def\pick{\textsf{pick}}
\def\int{\textsf{int}}
\def\x{\textsf{x}}
\def\y{\textsf{y}}
\def\aap{{\langle \auxb | \pubb \rangle}}
\def\aapc{{\langle \auxb | \pubb | \constb \rangle}}
\def\TreeRProofs{\textsf{TreeRProofs}}
\def\parentlabelsbits{\textsf{parent_labels_bits}}
\def\label{\textsf{label}}
\def\layerbits{\textsf{layer_bits}}
\def\labelbits{\textsf{label_bits}}
\def\digestbits{\textsf{digest_bits}}
\def\node{\textsf{node}}
\def\layerindex{\textsf{layer_index}}
\def\be{\textsf{be}}
\def\octet{\textsf{octet}}
\def\reverse{\textsf{reverse}}
\def\LSBit{\textsf{LSBit}}
\def\MSBit{\textsf{MSBit}}
\def\LSByte{\textsf{LSByte}}
\def\MSByte{\textsf{MSByte}}
\def\PorepPartitionProof{\textsf{PorepPartitionProof}}
\def\PostPartitionProof{\textsf{PostPartitionProof}}
\def\lebinrep#1{{\llcorner #1 \lrcorner_{\lower{2pt}{2, \textsf{le}}}}}
\def\bebinrep#1{{\llcorner #1 \lrcorner_{\lower{2pt}{2, \textsf{be}}}}}
\def\octetbinrep#1{{\llcorner #1 \lrcorner_{\lower{2pt}{2, \textsf{octet}}}}}
\def\fieldelement{\textsf{field_element}}
\def\Fqsafe{{\mathbb{F}_{q, \safe}}}
\def\elem{\textsf{elem}}
\def\challenge{\textsf{challenge}}
\def\challengeindex{\textsf{challenge_index}}
\def\uniquechallengeindex{\textsf{unique_challenge_index}}
\def\replicaindex{\textsf{replica_index}}
\def\uniquereplicaindex{\textsf{unique_replica_index}}
\def\nreplicas{\textsf{n_replicas}}
\def\unique{\textsf{unique}}
\def\R{\mathcal{R}}
\def\getpostchallenge{\textsf{get_post_challenge}}
\def\verifyvanillapostproof{\textsf{verify_vanilla_post_proof}}
\def\BinPathElement{\textsf{BinPathElement}}
\def\BinTreeDepth{\textsf{BinTreeDepth}}
\def\BinTree{\textsf{BinTree}}
\def\BinTreeProof{\textsf{BinTreeProof}}
\def\bintreeproofisvalid{\textsf{bintree_proof_is_valid}}
\def\Bit{{\{0, 1\}}}
\def\Byte{\mathbb{B}}
\def\calculatebintreechallenge{\textsf{calculate_bintree_challenge}}
\def\calculateocttreechallenge{\textsf{calculate_octtree_challenge}}
\def\depth{\textsf{depth}}
\def\dot{\textsf{.}}
\def\for{\textsf{for }}
\def\Function{\textbf{Function: }}
\def\Fq{{\mathbb{F}_q}}
\def\leaf{\textsf{leaf}}
\def\line#1#2#3{\scriptsize{\textsf{#1.}#2}\ \normalsize{#3}}
\def\missing{\textsf{missing}}
\def\NodeIndex{\textsf{NodeIndex}}
\def\nodes{\textsf{nodes}}
\def\OctPathElement{\textsf{OctPathElement}}
\def\OctTree{\textsf{OctTree}}
\def\OctTreeDepth{\textsf{OctTreeDepth}}
\def\OctTreeProof{\textsf{OctTreeProof}}
\def\octtreeproofisvalid{\textsf{octtree_proof_is_valid}}
\def\path{\textsf{path}}
\def\pathelem{\textsf{path_elem}}
\def\return{\textsf{return }}
\def\root{\textsf{root}}
\def\Safe{{\Byte^{[32]}_\textsf{safe}}}
\def\sibling{\textsf{sibling}}
\def\siblings{\textsf{siblings}}
\def\struct{\textsf{struct }}
\def\Teq{\underset{{\small \mathbb{T}}}{=}}
\def\Tequiv{\underset{{\small \mathbb{T}}}{\equiv}}
\def\thin{{\thinspace}}
\def\AND{\mathbin{\&}}
\def\MOD{\mathbin{\%}}
\def\createproof{\textsf{create_proof}}
\def\layer{\textsf{layer}}
\def\nodeindex{\textsf{node_index}}
\def\childindex{\textsf{child_index}}
\def\nodes{\textsf{nodes}}
\def\push{\textsf{push}}
\def\index{\textsf{index}}
\def\leaves{\textsf{leaves}}
\def\len{\textsf{len}}
\def\ColumnProof{\textsf{ColumnProof}}
\def\concat{\mathbin{\|}}
\def\inputs{\textsf{inputs}}
\def\Poseidon{\textsf{Poseidon}}
\def\bi{\ \ }
\def\Bool{{\{\textsf{True}, \textsf{False}\}}}
\def\curr{\textsf{curr}}
\def\if{\textsf{if }}
\def\else{\textsf{else}}
\def\proof{\textsf{proof}}
\def\Sha#1{\textsf{Sha#1}}
\def\ldotdot{{\ldotp\ldotp}}
\def\as{\textsf{ as }}
\def\bintreerootgadget{\textsf{bintree_root_gadget}}
\def\octtreerootgadget{\textsf{octtree_root_gadget}}
\def\cs{\textsf{cs}}
\def\RCS{\textsf{R1CS}}
\def\pathbits{\textsf{path_bits}}
\def\missingbit{\textsf{missing_bit}}
\def\missingbits{\textsf{missing_bits}}
\def\pubb{\textbf{pub}}
\def\privb{\textbf{priv}}
\def\auxb{\textbf{aux}}
\def\constb{\textbf{const}}
\def\CircuitVal{\textsf{CircuitVal}}
\def\CircuitBit{{\textsf{CircuitVal}_\Bit}}
\def\Le{\textsf{le}}
\def\privateinput{\textsf{private_input}}
\def\publicinput{\textsf{public_input}}
\def\deq{\mathbin{\overset{\diamond}{=}}}
\def\alloc{\textsf{alloc}}
\def\insertgadget#1{\textsf{insert_#1_gadget}}
\def\block{\textsf{block}}
\def\shagadget#1#2{\textsf{sha#1_#2_gadget}}
\def\poseidongadget#1{\textsf{poseidon_#1_gadget}}
\def\refeq{\mathbin{\overset{{\small \&}}=}}
\def\ptreq{\mathbin{\overset{{\small \&}}=}}
\def\bit{\textsf{bit}}
\def\extend{\textsf{extend}}
\def\auxle{{[\textbf{aux}, \textsf{le}]}}
\def\SpecificNotation{{\underline{\text{Specific Notation}}}}
\def\repeat{\textsf{repeat}}
\def\preimage{\textsf{preimage}}
\def\digest{\textsf{digest}}
\def\digestbytes{\textsf{digest_bytes}}
\def\digestint{\textsf{digest_int}}
\def\leencode{\textsf{le_encode}}
\def\ledecode{\textsf{le_decode}}
\def\ReplicaID{\textsf{ReplicaID}}
\def\replicaid{\textsf{replica_id}}
\def\replicaidbits{\textsf{replica_id_bits}}
\def\replicaidblock{\textsf{replica_id_block}}
\def\cc{\textsf{::}}
\def\new{\textsf{new}}
\def\lebitsgadget{\textsf{le_bits_gadget}}
\def\CircuitBitOrConst{{\textsf{CircuitValOrConst}_\Bit}}
\def\createporepcircuit{\textsf{create_porep_circuit}}
\def\CommD{\textsf{CommD}}
\def\CommC{\textsf{CommC}}
\def\CommR{\textsf{CommR}}
\def\CommCR{\textsf{CommCR}}
\def\commd{\textsf{comm_d}}
\def\commc{\textsf{comm_c}}
\def\commr{\textsf{comm_r}}
\def\commcr{\textsf{comm_cr}}
\def\assert{\textsf{assert}}
\def\asserteq{\textsf{assert_eq}}
\def\TreeDProof{\textsf{TreeDProof}}
\def\TreeRProof{\textsf{TreeRProof}}
\def\TreeR{\textsf{TreeR}}
\def\ParentColumnProofs{\textsf{ParentColumnProofs}}
\def\challengebits{\textsf{challenge_bits}}
\def\packedchallenge{\textsf{packed_challenge}}
\def\PartitionProof{\textsf{PartitionProof}}
\def\u#1{\textsf{u#1}}
\def\packbitsasinputgadget{\textsf{pack_bits_as_input_gadget}}
\def\treedleaf{\textsf{tree_d_leaf}}
\def\treerleaf{\textsf{tree_r_leaf}}
\def\calculatedtreedroot{\textsf{calculated_tree_d_root}}
\def\calculatedtreerleaf{\textsf{calculated_tree_r_leaf}}
\def\calculatedcommd{\textsf{calculated_comm_d}}
\def\calculatedcommc{\textsf{calculated_comm_c}}
\def\calculatedcommr{\textsf{calculated_comm_r}}
\def\calculatedcommcr{\textsf{calculated_comm_cr}}
\def\layers{\textsf{layers}}
\def\total{\textsf{total}}
\def\column{\textsf{column}}
\def\parentcolumns{\textsf{parent_columns}}
\def\columns{\textsf{columns}}
\def\parentlabel{\textsf{parent_label}}
\def\label{\textsf{label}}
\def\calculatedtreecleaf{\textsf{calculated_tree_c_leaf}}
\def\calculatedcolumn{\textsf{calculated_column}}
\def\parentlabels{\textsf{parent_labels}}
\def\drg{\textsf{drg}}
\def\exp{\textsf{exp}}
\def\parentlabelbits{\textsf{parent_label_bits}}
\def\parentlabelblock{\textsf{parent_label_block}}
\def\Bits{\textsf{ Bits}}
\def\safe{\textsf{safe}}
\def\calculatedlabel{\textsf{calculated_label}}
\def\createlabelgadget{\textsf{create_label_gadget}}
\def\encodingkey{\textsf{encoding_key}}
\def\encodegadget{\textsf{encode_gadget}}
\def\TreeC{\textsf{TreeC}}
\def\value{\textsf{value}}
\def\encoded{\textsf{encoded}}
\def\unencoded{\textsf{unencoded}}
\def\key{\textsf{key}}
\def\lc{\textsf{lc}}
\def\LC{\textsf{LC}}
\def\LinearCombination{\textsf{LinearCombination}}
\def\one{\textsf{one}}
\def\constraint{\textsf{constraint}}
\def\proofs{\textsf{proofs}}
\def\merkleproofs{\textsf{merkle_proofs}}
\def\TreeRProofs{\textsf{TreeRProofs}}
\def\challenges{\textsf{challenges}}
\def\pub{\textsf{pub}}
\def\priv{\textsf{priv}}
\def\last{\textsf{last}}
\def\TreeRProofs{\textsf{TreeRProofs}}
\def\post{\textsf{post}}
\def\SectorID{\textsf{SectorID}}
\def\winning{\textsf{winning}}
\def\window{\textsf{window}}
\def\Replicas{\textsf{Replicas}}
\def\P{\mathcal{P}}
\renewcommand\V{\mathcal{V}}
\def\ww{{\textsf{winning}|\textsf{window}}}
\def\replicasperk{{\textsf{replicas}/k}}
\def\replicas{\textsf{replicas}}
\def\Replica{\textsf{Replica}}
\def\createvanillapostproof{\textsf{create_vanilla_post_proof}}
\def\createpostcircuit{\textsf{create_post_circuit}}
\def\ReplicaProof{\textsf{ReplicaProof}}
\def\aww{{\langle \ww \rangle}}
\def\partitionproof{\textsf{partition_proof}}
\def\replicas{\textsf{replicas}}
\def\getdrgparents{\textsf{get_drg_parents}}
\def\getexpparents{\textsf{get_exp_parents}}
\def\DrgSeed{\textsf{DrgSeed}}
\def\DrgSeedPrefix{\textsf{DrgSeedPrefix}}
\def\FeistelKeysBytes{\textsf{FeistelKeysBytes}}
\def\porep{\textsf{porep}}
\def\rng{\textsf{rng}}
\def\ChaCha#1{\textsf{ChaCha#1}}
\def\cc{\textsf{::}}
\def\fromseed{\textsf{from_seed}}
\def\buckets{\textsf{buckets}}
\def\meta{\textsf{meta}}
\def\dist{\textsf{dist}}
\def\each{\textsf{each}}
\def\PorepID{\textsf{PorepID}}
\def\porepgraphseed{\textsf{porep_graph_seed}}
\def\utf{\textsf{utf8}}
\def\DrgStringID{\textsf{DrgStringID}}
\def\FeistelStringID{\textsf{FeistelStringID}}
\def\graphid{\textsf{graph_id}}
\def\createfeistelkeys{\textsf{create_feistel_keys}}
\def\FeistelKeys{\textsf{FeistelKeys}}
\def\feistelrounds{\textsf{fesitel_rounds}}
\def\feistel{\textsf{feistel}}
\def\ExpEdgeIndex{\textsf{ExpEdgeIndex}}
\def\loop{\textsf{loop}}
\def\right{\textsf{right}}
\def\left{\textsf{left}}
\def\mask{\textsf{mask}}
\def\RightMask{\textsf{RightMask}}
\def\LeftMask{\textsf{LeftMask}}
\def\roundkey{\textsf{round_key}}
\def\beencode{\textsf{be_encode}}
\def\Blake{\textsf{Blake2b}}
\def\input{\textsf{input}}
\def\output{\textsf{output}}
\def\while{\textsf{while }}
\def\digestright{\textsf{digest_right}}
\def\xor{\mathbin{\oplus_\text{xor}}}
\def\Edges{\textsf{ Edges}}
\def\edge{\textsf{edge}}
\def\expedge{\textsf{exp_edge}}
\def\expedges{\textsf{exp_edges}}
\def\createlabel{\textsf{create_label}}
\def\Label{\textsf{Label}}
\def\Column{\textsf{Column}}
\def\Columns{\textsf{Columns}}
\def\ParentColumns{\textsf{ParentColumns}}
\def\tern#1?#2:#3{#1\ \text{?}\ #2 \ \text{:}\ #3}
\def\repeattolength{\textsf{repeat_to_length}}
\def\verifyvanillaporepproof{\textsf{verify_vanilla_porep_proof}}
\def\poreppartitions{\textsf{porep_partitions}}
\def\challengeindex{\textsf{challenge_index}}
\def\porepbatch{\textsf{porep_batch}}
\def\winningchallenges{\textsf{winning_challenges}}
\def\windowchallenges{\textsf{window_challenges}}
\def\PorepPartitionProof{\textsf{PorepPartitionProof}}
\def\TreeD{\textsf{TreeD}}
\def\TreeCProof{\textsf{TreeCProof}}
\def\Labels{\textsf{Labels}}
\def\porepchallenges{\textsf{porep_challenges}}
\def\postchallenges{\textsf{post_challenges}}
\def\PorepChallengeSeed{\textsf{PorepChallengeSeed}}
\def\getporepchallenges{\textsf{get_porep_challenges}}
\def\getallparents{\textsf{get_all_parents}}
\def\PorepChallengeProof{\textsf{PorepChallengeProof}}
\def\challengeproof{\textsf{challenge_proof}}
\def\PorepChallenges{\textsf{PorepChallenges}}
\def\replicate{\textsf{replicate}}
\def\createreplicaid{\textsf{create_replica_id}}
\def\ProverID{\textsf{ProverID}}
\def\replicaid{\textsf{replica_id}}
\def\generatelabels{\textsf{generate_labels}}
\def\labelwithdrgparents{\textsf{label_with_drg_parents}}
\def\labelwithallparents{\textsf{label_with_all_parents}}
\def\createtreecfromlabels{\textsf{create_tree_c_from_labels}}
\def\ColumnDigest{\textsf{ColumnDigest}}
\def\encode{\textsf{encode}}
$$

# SDR Spec

## Merkle Proofs

**Implementation:**
* [`storage_proofs::merkle::MerkleTreeWrapper::gen_proof()`]()
* [`merkle_light::merkle::MerkleTree::gen_proof()`](https://github.com/filecoin-project/merkle_light/blob/64a468807c594d306d12d943dd90cc5f88d0d6b0/src/merkle.rs#L918)

**Additional Notation:**
$\index_l: [\lfloor N_\nodes / 2^l \rfloor] \equiv [\len(\BinTree\dot\layer_l)]$
The index of a node in a $\BinTree$ layer $l$. The leftmost node in a tree has $\index_l = 0$. For each tree layer $l$ (excluding the root layer) a Merkle proof verifier calculates the label of the node at $\index_l$ from a single Merkle proof path element $\BinTreeProof_c\dot\path[l - 1] \thin$.

### BinTreeProofs

The method $\BinTreeProof\dot\createproof$ is used to create a Merkle proof for a challenge node $c$.

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

**Code Comments:**
* **Line 5:** Calculates the node index in layer $l$ of the node that the verifier calculated using the previous lath element (or the $\BinTreeProof_c\dot\leaf$ if $l = 0). Note that $c \gg l \equiv \lfloor c / 2^l \rfloor \thin$.

### OctTreeProofs

The method $\OctTreeProof\dot\createproof$ is used to create a Merkle proof for a challenge node $c$.

**Additional Notation:**
$\index_l: [\lfloor N_\nodes / 8^l \rfloor] \equiv [\len(\OctTree\dot\layer_l)]$
The index of a node in an $\OctTree$ layer $l$. The leftmost node in a tree has $\index_l = 0$. For each tree layer $l$ (excluding the root layer) a Merkle proof verifier calculates the label of the node at $\index_l$ from a single Merkle proof path element $\OctTreeProof_c\dot\path[l - 1] \thin$.

$\textsf{first_sibling}_l \thin, \textsf{last_sibling}: [\lfloor N_\nodes / 8^l \rfloor]$
The node indexes in tree layer $l$ of the first and last nodes in this layer's Merkle path element's siblings array $\OctTreeProof_c\dot\path[l]\dot\siblings \thin$.

$\overline{\underline{\Function \OctTree\dot\createproof(c: \NodeIndex) \rightarrow \OctTreeProof_c}}$
$\line{1}{\bi}{\leaf: \Fq = \OctTree\dot\leaves[c]}$
$\line{2}{\bi}{\root: \Fq = \OctTree\dot\root}$

$\line{3}{\bi}{\path: \OctPathElement^{[\OctTreeDepth]}= [\ ]}$
$\line{4}{\bi}{\for l \in [\OctTreeDepth]:}$
$\line{5}{\bi}{\quad \index_l: [\len(\OctTree\dot\layer_l)] = c \gg (3 * l)}$
$\line{6}{\bi}{\quad \missing: [8] = \index_l \MOD 8}$

$\line{7}{\bi}{\quad \textsf{first_sibling}_l = \index_l - \missing}$
$\line{8}{\bi}{\quad \textsf{last_sibling}_l = \index_l + (7 - \missing)}$
$\line{9}{\bi}{\quad \siblings: \Fq^{[7]} =}$
$\quad\quad\quad \OctTree\dot\layer_l[\textsf{first_sibling}_l \mathbin{\ldotdot} \index_l]$
$\quad\quad\quad \|\ \OctTree\dot\layer_l[\index_l + 1 \mathbin{\ldotdot} \textsf{last_sibling}_l + 1]$

$\line{10}{}{\quad \path\dot\push(\OctPathElement \thin \{\ \siblings, \thin \missing\ \} \thin )}$
$\line{11}{}{\return \OctTreeProof_c \thin \{\ \leaf, \thin \root, \thin \path\ \}}$

**Code Comments:**
* **Line 5:** Calculates the node index in layer $l$ of the node that the verifier calculated themselves using the previous path element (or $\OctTreeProof_c\dot\leaf$ if $l = 0$). Note that $c \gg (3 * l) \equiv \lfloor c / 8^l \rfloor \thin$.
* **Line 7-8:** Calculates the indexes in tree layer $l$ of the first and last (both inclusive) Merkle hash inputs for layer $l$'s path element.
* **Line 9:** Copies the 7 Merkle hash inputs that will be in layer $l$'s path element $\OctTreeProof\dot\path[l]\dot\siblings \thin$.

### Proof Root Validation

The functions $\bintreeproofisvalid$ and $\octtreeproofisvalid$ are used to verify that a $\BinTreeProof\dot\path$ or an $\OctTreeProof\dot\path$ hash to the root found in the Merkle proof $\BinTreeProof\dot\root$ and $\OctTreeProof\dot\root$ respectively. 

Note that these functions do not verify that a $\BinTreeProof\dot\path$ or an $\OctTreeProof\dot\path$ correspond to the expected Merkle challenge $c$. To verify that a proof path is consistent with $c$, see the psuedocode functions $\calculatebintreechallenge$ and $\calculateocttreechallenge$.

**Implementation:**
* [`storage_proofs::core::merkle::proof::SingleProof::verify()`](https://github.com/filecoin-project/rust-fil-proofs/blob/f75804c503d9b97a2b02ef3ea2e5d44e8e2c6470/storage-proofs/core/src/merkle/proof.rs#L540)
* [`storage_proofs::core::merkle::proof::InclusionPath::root()`](https://github.com/filecoin-project/rust-fil-proofs/blob/f75804c503d9b97a2b02ef3ea2e5d44e8e2c6470/storage-proofs/core/src/merkle/proof.rs#L189)

$\overline{\underline{\Function \bintreeproofisvalid(\proof: \BinTreeProof) \rightarrow \Bool}\thin}$
$\line{1}{\bi}{\curr: \Safe = \proof\dot\leaf}$
$\line{2}{\bi}{\for \sibling, \missing \in \proof\dot\path:}$
$\line{3}{\bi}{\quad \curr: \Safe = \if \missing = 0:}$
$\quad\quad\quad\quad \Sha{254}_2([\curr, \sibling])$
$\quad\quad\thin \else:$
$\quad\quad\quad\quad \Sha{254}_2([\sibling, \curr])$
$\line{4}{\bi}{\return \curr = \proof\dot\root}$

The function $\octtreeproofisvalid$ can receive as the type of its $\proof$ argument either an $\OctTreeProof$ or $\ColumnProof$ (a $\ColumnProof$ is just an $\OctTreeProof$ with an adjoined field $\column$, $\ColumnProof_c \equiv \OctTreeProof_c \cup \column_c \thin$).

$\overline{\underline{\Function \octtreeproofisvalid(\proof: \OctTreeProof) \rightarrow \Bool}\thin}$
$\line{1}{\bi}{\curr: \Fq = \proof\dot\leaf}$
$\line{2}{\bi}{\for \siblings, \missing \in \proof\dot\path:}$
$\line{3}{\bi}{\quad \inputs: \Fq^{[8]} = \siblings[\ldotdot \missing] \concat \curr \concat \siblings[\missing \ldotdot]}$
$\line{4}{\bi}{\quad \curr = \Poseidon_8(\inputs)}$
$\line{5}{\bi}{\return \curr = \proof\dot\root}$

$\overline{\underline{\Function \octtreeproofisvalid(\proof: \ColumnProof) \rightarrow \Bool}\thin}$
$\line{1}{\bi}{\return \octtreeproofisvalid(\OctTreeProof\ \{\ \leaf, \root, \path \Leftarrow \ColumnProof\ \})}$

### Merkle Proof Challenge Validation

Given a Merkle path $\path$ in a $\BinTree$ or $\OctTree$, $\calculatebintreechallenge$ and $\calculateocttreechallenge$ calculate the Merkle challenge $c$ for which the Merkle proof $\path$ was generated.

Given a Merkle challenge $c$ and its path in a $\BinTree$ or $\OctTree$, the concatentation of the $\missing$ bits (or octal digits) in the Merkle path is the little-endian binary (or octal) representation of the integer $c \thin$:

$\quad \line{1}{\bi}{c: \NodeIndex = \langle \text{challenge} \rangle}$
$\quad \line{2}{\bi}{\BinTreeProof_c = \BinTree\dot\createproof(c)}$
$\quad \line{3}{\bi}{\OctTreeProof_c = \OctTree\dot\createproof(c)}$
$\quad \line{4}{\bi}{\llcorner c \lrcorner_{2, \Le}: \Bit^{[\hspace{1pt} \log_2(N_\nodes) \hspace{1pt}]} = \big\|_{\pathelem \hspace{1pt} \in \hspace{1pt} \BinTreeProof_c\dot\path} \thin \pathelem\dot\missing} \vphantom{{{|^|}^|}^x}$
$\quad \line{5}{\bi}{\rlap{\llcorner c \lrcorner_{8, \Le}: [8]^{[\hspace{1pt} \log_8(N_\nodes) \hspace{1pt}]}}\hphantom{\llcorner c \lrcorner_{2, \Le}: \Bit^{[\hspace{1pt} \log_2(N_\nodes) \hspace{1pt}]}} = \big\|_{\pathelem \hspace{1pt} \in \hspace{1pt} \BinTreeProof_c\dot\path} \thin \pathelem\dot\missing} \vphantom{{{|^|}^|}^x}$
$\quad \line{6}{\bi}{\rlap{\llcorner c \lrcorner_{2, \Le}: \Bit^{[\hspace{1pt} \log_2(N_\nodes) \hspace{1pt}]}}\hphantom{\llcorner c \lrcorner_{2, \Le}: \Bit^{[\hspace{1pt} \log_2(N_\nodes) \hspace{1pt}]}} = \big\|_{\pathelem \hspace{1pt} \in \hspace{1pt} \OctTreeProof_c\dot\path} \thin \llcorner \pathelem\dot\missing \lrcorner_{2, \Le}} \vphantom{{|^|}^x}$

**Implementation:** [`storage_proofs::merkle::MerkleProofTrait::path_index()`](https://github.com/filecoin-project/rust-fil-proofs/blob/f75804c503d9b97a2b02ef3ea2e5d44e8e2c6470/storage-proofs/core/src/merkle/proof.rs#L89)

**Additional Notation:**
$\path = \BinTreeProof\dot\path$
$\path = \OctTreeProof\dot\path$
The $\path$ argument is the path field of a $\BinTreeProof$ or $\OctTreeProof$.

$c: \NodeIndex$
The challenge corresponding to $\path$.

$l \in [\BinTreeDepth]$
$l \in [\OctTreeDepth]$
A path element's layer in a Merkle tree (the layer in the tree that contains the path elements $\siblings$). Layer $l = 0$ is the leaves layer of the tree. Here, values for $l$ do not include the root layer $l \neq \BinTreeDepth, \OctTreeDepth \thin$.

$\overline{\underline{\Function \calculatebintreechallenge(\path: \BinPathElement^{[\BinTreeDepth]}) \rightarrow c}}$
$\line{1}{\bi}{\return \sum_{l \in [\BinTreeDepth]}{\path[l]\dot\missing * 2^l}}$

$\overline{\underline{\Function \calculateocttreechallenge(\path: \OctPathElement^{[\OctTreeDepth]}) \rightarrow c}}$
$\line{1}{\bi}{\return \sum_{l \in [\OctTreeDepth]}{\path[l]\dot\missing * 8^l}}$

## Stacked Depth Robust Graphs

Filecoin utilizes the topological properties of depth robust graphs (DRG's) to build a sequential and regeneration resistant encoding scheme. We stack $N_\layers$ of DRG's, each containing $N_\nodes$ nodes, on top of one another and connect each adjacent pair of DRG layers via the edges of bipartite expander. The source layer of each expander is the DRG at layer $l$ and the sink layer is the DRG at layer $l + 1$. The resulting graph is termed a Stacked-DRG. 

For every node $v \in [N_\nodes]$ in the DRG at layer $l \in [N_\layers]$, we generate $d_\drg$ number of DRG parent for $v$ in layer $l$. DRG parents are generated using the [Bucket Sampling algorithm](https://eprint.iacr.org/2017/443.pdf). For every node $v$ in layers $l > 0$ we generate $d_\exp$ number of expander parents for $v$ where the parents are in layer $l - 1$. Expander parents are generated using a psuedorandom permutation (PRP) $\pi: [N_\nodes] \rightarrow [N_\nodes]$ which maps a node in the DRG layer $l$ to a node in the the DRG layer $l - 1$. The psudeorandom permutation is generated using a $N_\fesitelrounds$-round Feistel network whose round function is the keyed hash function $\Blake$, where round keys are specified by the constant $\FeistelKeys$.

### DRG

The function $\getdrgparents$ are used to generate a node $v$'s DRG parents in the Stacked-DRG layer $l_v \thin$. The set of DRG parents returned for $v$ is the same for all PoRep's of the same PoRep version $\PorepID$.

**Implementation:** [`storage_proofs::core::drgraph::BucketGraph::parents()`](https://github.com/filecoin-project/rust-fil-proofs/blob/d15b4307abe384d49ab2ce76b57377204f935cec/storage-proofs/core/src/drgraph.rs#L130)

**Additional Notation:**
$v, u: \NodeIndex$
DRG child and parent node indexes respectively. A DRG child and its parents are in the same Stacked-DRG layer $l$.

$\mathbf{u}_\drg: \NodeIndex^{[d_\drg]}$
The set of $v$'s DRG parents.

$v_\meta, u_\meta: [d_\meta * N_\nodes]$
The indexes of $v$ and $u$ in the metagraph.

$\rng_{\PorepID, v}$
The RNG used to sample $v$'s DRG parents. The RNG is seeded with the same bytes $\DrgSeed_{\PorepID, v}$ every time $\getdrgparents$ is called for node $v$ from a PoRep having version $\PorepID$.

$x \xleftarrow[\rng]{} S$
Samples $x$ uniformly from $S$ using the seeded $\rng$.

$b: [1, N_\buckets + 1]$
The Bucket Sampling bucket index. **Bucket indexes start at 1.**

$\dist_{\min, b}$
$\dist_{\max, b}$
The minimum and maximum parent distances in bucket $b$.

$\dist_{u_\meta}$
The distance $u_\meta$ is from $v_\meta$ in the metagraph.

$\overline{\underline{\Function \getdrgparents(v: \NodeIndex) \rightarrow \NodeIndex^{[d_\drg]}}}$
$\line{1}{\bi}{\DrgSeed_{\PorepID, v}: \Byte^{[32]} = \DrgSeed_\PorepID \concat \leencode(v) \as \Byte^{[4]}}$
$\line{2}{\bi}{\rng_{\PorepID, v} = \ChaCha{8}\cc\fromseed(\DrgSeed_{\PorepID, v})}$

$\line{3}{\bi}{\mathbf{u}_\drg: \textsf{NodeIndex}^{[d_\drg]} = [\ ]}$

$\line{4}{\bi}{v_\meta = v * d_\textsf{meta}}$
$\line{5}{\bi}{N_\buckets = \lceil \log_2(v_\meta) \rceil}$
$\line{6}{\bi}{\for \each \in [d_\meta]:}$
$\line{7}{\bi}{\quad b \xleftarrow[\rng]{}  [1, N_\buckets + 1]}$
$\line{8}{\bi}{\quad \dist_{\max, b} = \textsf{min}(v_\meta, 2^b)}$
$\line{9}{\bi}{\quad \dist_{\min, b} = \textsf{max}(d_{\max, b} / 2, 2)}$
$\line{10}{}{\quad \dist_{u_\meta} \xleftarrow[\rng]{} [\dist_{\min, b} \thin, \dist_{\max, b}]}$
$\line{11}{}{\quad u_\meta = v_\meta - \dist_{u_\meta}}$
$\line{12}{}{\quad u: \NodeIndex = \lfloor u_\meta / d_\meta \rfloor}$
$\line{13}{}{\quad \mathbf{u}_\drg\dot\push(u)}$

$\line{14}{}{\mathbf{u}_\drg\dot\push(v - 1)}$
$\line{15}{}{\return \mathbf{u}_\drg}$

### Expander

The function $\getexpparents$ is used to generate a node $v$'s expander parents in the Stacked-DRG layer $l_v - 1 \thin$. The set of expander parents returned for a node $v$ is the same for all PoRep's of the same version $\PorepID$.

**Implementation:** [`storage_proofs::porep::stacked::vanilla::graph::StackedGraph::generate_expanded_parents()`](https://github.com/filecoin-project/rust-fil-proofs/blob/d15b4307abe384d49ab2ce76b57377204f935cec/storage-proofs/porep/src/stacked/vanilla/graph.rs#L448)

**Additional Notation:**
$v, u: \NodeIndex$
Expander child and parent node indexes respectively. Each expander parent $u$ is in the Staked-DRG layer $l - 1$ that precedes the child node $v$'s layer $l$.

$\mathbf{u}_\exp: \NodeIndex^{[d_\exp]}$
The set of $v$'s expander parents.

$e_l \thin, e_{l - 1}: \ExpEdgeIndex$
The index of an expander edge in the child $v$'s layer $l$ and the parent $u$'s layer $l - 1$ respectively. An expander edge connects edge indexes $(e_{l - 1}, e_l)$ in adjacent Stacked-DRG layers.

$\overline{\underline{\Function \getexpparents(v: \NodeIndex) \rightarrow \NodeIndex^{[d_\exp]}}}$
$\line{1}{\bi}{\mathbf{u}_\exp: \NodeIndex^{[d_\exp]} = [\ ]}$
$\line{2}{\bi}{\for p\in [d_E]:}$
$\line{3}{\bi}{\quad e_l = v * d_E + p}$
$\line{4}{\bi}{\quad e_{l - 1} = \feistel(e_l)}$
$\line{5}{\bi}{\quad u: \NodeIndex = \lfloor e_{l - 1} / d_\exp \rfloor}$
$\line{6}{\bi}{\quad \mathbf{u}_\exp\dot\push(u)}$
$\line{7}{\bi}{\return \mathbf{u}_\exp}$

### Feistel Network PRP

The function $\feistel$ runs an $N_\feistelrounds = 4$ round Feistel network as a psuedorandom permutation (PRP) over the set of expander edges $[d_\exp * N_\nodes] = [2^{33}]$ in a Stacked-DRG layer.

**Implementation:** [`storage_proofs::core::crypto::feistel::permute()`](https://github.com/filecoin-project/rust-fil-proofs/blob/d15b4307abe384d49ab2ce76b57377204f935cec/storage-proofs/core/src/crypto/feistel.rs#L34)

**Additional Notation:**
$\input, \output$
The Feistel network's input and output blocks respectively. The $\input$ argument and the returned $\output$ are guaranteed to be valid $\u{33}$ expander edge indexes, however their intermediate values may be $\u{34}$.

$\u{64}_{(17)}$
An unsigned 64-bit integer whose lowest 17 bits are utilized. The integer's 17 least significant bits are used and may be $0$ or $1$, while the integer's 47 most significant bits are $0$ and are not used.

$\left_r, \right_r$
The left and right halves of round $r$'s input block.

$\left_{r + 1}, \right_{r + 1}$
The left and right halves of the next round's $r + 1$ input block.

$\FeistelKeys_\PorepID$
Is the set of constant round keys associated with the PoRep version $\PorepID$ that called $\feistel$.

$\key_r$
Round $r$'s key.

$\digestright$
The $\ell_\mask^\bit = 17$ least significant bits of a round's Blake2b $\digest$.

$\overline{\underline{\Function \feistel(\input: \ExpEdgeIndex) \rightarrow \ExpEdgeIndex\ }\thin}$
$\line{1}{\bi}{\output: \u{64}_{(34)} = \varnothing}$

$\line{2}{\bi}{\while \output \notin [N_\expedges]:}$
$\line{3}{\bi}{\quad \right_r: \u{64}_{(17)} = \input \AND \mathsf{RightMask}}$
$\line{4}{\bi}{\quad \left_r: \u{64}_{(17)} = (\input \AND \LeftMask) \gg \ell_\mask^\bit}$

$\line{5}{\bi}{\quad \for \key_r \in \FeistelKeys_\PorepID:}$
$\line{6}{\bi}{\quad\quad \preimage: \Byte^{[16]} = \beencode(\right_r) \as \Byte^{[8]} \concat \beencode(\key_r) \as \Byte^{[8]}}$
$\line{7}{\bi}{\quad\quad \digest: \Byte^{[8]} = \Blake(\preimage)[..8]}$
$\line{8}{\bi}{\quad\quad \digest: \u{64} = \bedecode(\digest)}$
$\line{9}{\bi}{\quad\quad \digestright: \u{64}_{(17)} = \digest \AND \RightMask}$

$\line{10}{}{\quad\quad \left_{r + 1}: \u{64}_{(17)} = \right_r}$
$\line{11}{}{\quad\quad \right_{r + 1}: \u{64}_{(17)} = \left_r \xor \digestright}$

$\line{12}{}{\quad\quad \left_r, \right_r = \left_{r + 1}, \right_{r + 1}}$

$\line{13}{}{\quad \output: \u{64}_{(34)} = (\left_r \ll \ell_\mask^\bit) \OR \right_r}$
$\line{14}{}{\return \output}$

**Code Comments:**
* **Line 1:** Initializes the output of the network to the null value $\varnothing \thin$. The value of $\output$ is a $34 = 2 * \ell_\mask^\bit$-bit integer and may not be a valid expander graph edge index after $N_\feistelrounds$ rounds.
* **Line 2:** The network runs at least once and reruns if $\output$'s most significant bit (its 34th bit) is set to $1$ ($\output$ is not a valid $\u{33}$ expander edge index).

### All Parents

The function $\getallparents$ returns a node $v$'s set of DRG parents concatenated with its set of expander parents. The set of parents returned for $v$ is the same across all PoRep's of the same PoRep version $\PorepID$.

$\overline{\underline{\Function \getallparents(v: \NodeIndex) \rightarrow \mathbf{u}_\total}}$
$\line{1}{\bi}{\return \getdrgparents(v) \concat \getexpparents(v)}$

## Labeling

### Labeling a Node

The labeling function for every node in a Stacked-DRG is $\Sha{254}$ producing a 254-bit field element $\Fqsafe$. A unique preimage is derived for each node-layer tuple $(v, l)$ in replica $\ReplicaID$'s Stacked-DRG.

The labeling preimage for the first node $v_0 = 0$ in every Stacked-DRG layer $l \in [N_\layers]$ for a replica $\ReplicaID$ is defined:

$\bi \preimage_{v_0, l}: \Byte^{[44]} =$
$\bi\quad\quad \ReplicaID \concat \beencode(l \as \u{32}) \as \Byte^{[4]} \concat \beencode(v_0 \as \u{64}) \as \Byte^{[8]}$

The labeling preimage for each node $v > 0$ in the first layer $l_0 = 0$ is defined:

$\bi \preimage_{v, l_0}: \Byte^{[1228]} =$
$\bi\quad\quad \ReplicaID$
$\bi\quad\quad \|\ \beencode(l_0 \as \u{32}) \as \Byte^{[4]}$
$\bi\quad\quad \|\ \beencode(v \as \u{64}) \as \Byte^{[8]}$
$\bi\quad\quad \big{\|}_{\Label_{u, l_0} \hspace{1pt} \in \hspace{1pt} \ParentLabels_{\mathbf{u}_\drg}^\star} \Label_{u, l_0} \as \Byte^{[32]} \vphantom{|^|}$

The labeling preimage for each node $v > 0$ in each layer $l > 0$ is defined:

$\bi \preimage_{v, l}: \Byte^{[1228]} =$
$\bi\quad\quad \ReplicaID$
$\bi\quad\quad \|\ \beencode(l \as \u{32}) \as \Byte^{[4]}$
$\bi\quad\quad \|\ \beencode(v \as \u{64}) \as \Byte^{[8]}$
$\bi\quad\quad \big{\|}_{\Label_u \hspace{1pt} \in \hspace{1pt} \ParentLabels_{\mathbf{u}_\total}^\star} \Label_u \as \Byte^{[32]} \vphantom{|^|}$

### The Labels Matrix

The function $\textsf{generate_labels}$ describes how every Stacked-DRG node is labeled for a replica. Nodes in the first layer $l_0 = 0$ are labeled using only DRG parents' labels, nodes in every subsequent layers $l > 0$ are labeled using both their DRG and expander parents' labels. The first node $v_0$ in every layer is not labeled using parents.

**Additional Notation:**
$\Labels_R$
Denotes that $\Labels$ is the labeling for a single replica $R$.

$l_0 = 0$
The constant $l_0$ is used to signify the first Stacked-DRG layer.

$l_v$
The Stacked-DRG layer in which a node $v$ resides.

$\Label_{v, l_v}$
The label of node $v$ in the Stacked-DRG layer $l_v$.

$u_{\langle \drg | \exp \rangle}$
Denotes that parent $u$ may be a DRG or expander parent. 

$\Label_{u_\drg}$
The label of a DRG parent (in $v$'s layer $l$).

$\Label_u \equiv \or{\Label_{u_\drg, l_v}}{\Label_{u_\exp, l_v - 1}}$
The label of either a DRG or expander parent (in layer $l$ or $l - 1$ respectively).

$\overline{\underline{\Function \generatelabels(\ReplicaID) \rightarrow \Labels_R}}$
$\line{1}{\bi}{\Labels: {\Label^{[N_\nodes]}}^{[N_\layers]} 
= [\ ]}$

$\line{2}{\bi}{\for v \in [N_\nodes]:}$
$\line{3}{\bi}{\quad \Labels[l_0][v] = \labelwithdrgparents(\ReplicaID, v, \Labels[l_0][..v])}$

$\line{4}{\bi}{\for l \in [1, N_\layers - 1]:}$
$\line{5}{\bi}{\quad \for v \in [N_\mathsf{nodes}]:}$
$\line{6}{\bi}{\quad\quad \Labels[l][v] = \labelwithallparents(\ReplicaID, l, v, \Labels[l][..v], \Labels[l - 1])}$

$\line{7}{\bi}{\return \Labels}$

**Code Comments:**
* **Lines 2-3:** Label the first Stacked-DRG layer.
* **Lines 4-6:** Label the remaining Stacked-DRG layers.

<br />

The function $\labelwithdrgparents$ is used to label a node $v$ in the first Stacked-DRG layer $l_0 = 0 \thin$.

The label of each node $v$ in $l_0$ is dependent on the labels of $v$'s DRG parents (where $v$'s DRG parents are in layer $l_v = l_0 \thin$. $v$'s DRG parents always come from the node range $[v]$ in layer $l_v$, thus we pass in the argument $\Labels[l_0][..v]$ which contains the label of every node up to $v$ in $l_0 \thin$. $\Labels[l_0][..v]$ is guaranteed to be labeled up to node $v$ because $\labelwithdrgparents$ is called sequentially for each node $v \in [N_\nodes] \thin$.

$\overline{\underline{\Function \labelwithdrgparents(\ReplicaID, v: \NodeIndex, \Labels[l_0][..v]) \rightarrow \Label_{v, l_0}}}$
$\line{1}{\bi}{\preimage: \Byte^{[*]} =}$
$\quad\quad \ReplicaID \concat \beencode(l_0 \as \u{32}) \as \Byte^{[4]} \concat \beencode(v \as \u{64}) \as \Byte^{[8]}$

$\line{2}{\bi}{\if v > 0:}$
$\line{3}{\bi}{\quad \mathbf{u}_\drg: \textsf{NodeIndex}^{[d_\drg]} = \getdrgparents(v)}$
$\line{4}{\bi}{\quad \for i \in [N_\parentlabels]:}$
$\line{5}{\bi}{\quad\quad u_\drg = \mathbf{u}_\drg[i \MOD d_\drg]}$
$\line{6}{\bi}{\quad\quad \Label_{u_\drg, l_0}: \Fqsafe = \Labels[l_0][u_\drg]}$
$\line{7}{\bi}{\quad\quad \preimage\dot\extend(\leencode(\Label_{u_\drg, l_0}) 
\as \Safe)}$

$\line{8}{\bi}{\return \Sha{254}(\preimage) \as \Fqsafe}$

<br />

The function $\labelwithallparents$ is used to label a node $v$ in all layers other than the first Stacked-DRG layer $l_v > 0 \thin$.

The label of a node $v$ in layers $l_v > 0$ is dependent on both the labels of $v$'s DRG and expander parents. $\labelwithallparents$ takes the argument $\Labels[l_v][..v]$ (the current layer $l_v$ being labeled, contains labels up to node $v$) to retrieve the labels of $v$'s DRG parents and the argument $\Labels[l_v - 1]$ (the previous layer's labels) to retrieve the labels of $v$'s expander parents.

$\overline{\Function \labelwithallparents(\bi}$
$\quad \ReplicaID,$
$\quad l_v \in [1, N_\layers - 1],$
$\quad v: \NodeIndex,$
$\quad \Labels[l_v][..v],$
$\quad \Labels[l_v - 1],$
$\underline{) \rightarrow \Label_{v, l_v} \qquad\qquad\qquad\qquad\qquad\bi}$
$\line{1}{\bi}{\preimage: \Byte^{[*]} =}$
$\quad\quad \ReplicaID \concat \beencode(l_v \as \u{32}) \as \Byte^{[4]} \concat \beencode(v \as \u{64}) \as \Byte^{[8]}$

$\line{2}{\bi}{\if v > 0:}$
$\line{3}{\bi}{\quad \mathbf{u}_\total: \textsf{NodeIndex}^{[d_\total]} = \getallparents(v)}$
$\line{4}{\bi}{\quad \for i \in [N_\parentlabels]:}$
$\line{5}{\bi}{\quad\quad p = i \MOD d_\total}$
$\line{6}{\bi}{\quad\quad u_{\langle \drg | \exp \rangle} = \mathbf{u}_\total[p]}$
$\line{7}{\bi}{\quad\quad \Label_u: \Fqsafe = \if p < d_\drg}:$
$\quad\quad\quad\quad \Labels[l_v][u_\drg]$
$\quad\quad\quad \else:$
$\quad\quad\quad\quad \Labels[l_v - 1][u_\exp]$
$\line{8}{\bi}{\quad\quad \preimage\dot\extend(\leencode(\Label_u) 
\as \Safe)}$

$\line{9}{\bi}{\return \Sha{254}(\preimage) \as \Fqsafe}$

## Column Commitments

The column commitment process is used commit to a replica's labeling $\Labels$. The column commitment $\CommC$ is generated by building an $\TreeC: \OctTree$ over the labeling matrix $\Labels$ and taking the tree's root.

To build a tree over the matrix $\Labels$ we hash each of its $N_\nodes$ number of columns (where each column contains $N_\layers$ number of $\Label$'s) using the hash function $\Poseidon_{11}$ producing $N_\nodes$ number of column digests. The $i^{th}$ column digest is the $i^{th}$ leaf in $\TreeC$.

$\overline{\underline{\Function \createtreecfromlabels(\Labels) \rightarrow \TreeC}}$
$\line{1}{\bi}{\leaves: {\Fq}^{[N_\nodes]} = [\ ]}$
$\line{2}{\bi}{\for v \in [N_\nodes]:}$
$\line{3}{\bi}{\quad \column_v: \Fqsafe^{[N_\layers]} = \Labels[:][v]}$
$\line{4}{\bi}{\quad \digest: \Fq = \Poseidon_{11}(\column_v)}$
$\line{5}{\bi}{\quad \leaves\dot\push(\digest)}$
$\line{6}{\bi}{\return \OctTree\cc\new(\leaves)}$

## Encoding

Encoding is the process by which a sector $D: \Safe^{[N_\nodes]}$ is transformed into its encoding $R: \Fqsafe^{[N_\nodes]}$. The encoding function is *node-wise* prime field addition $\oplus$, where "node-wise" means that every distinct $\Safe$ slice $D_i \in D$ is discretely encoded.

$D$ is viewed as an array of $N_\nodes$ distinct byte arrays $D_i: \Safe$. Sector preprocessing ensures that each $D_i$ is a valid $\Safe$ (represents a valid 254-bit or less field element $\Fqsafe)$.

$\bi D: \Safe^{[N_\nodes]} = [D_0, \ldots, D_{N_\nodes - 1}]$
$\bi D_i: \Safe = D[i * 32 \thin\ldotdot\thin (i + 1) * 32]$

A unique encoding key $K$ is derived for every distinct $\ReplicaID$ via the PoRep labeling process producing $\Labels$. Each $D_i \in D$ is encoded by a distinct encoding key $K_i \in K$, where $K_i$ is $i^{th}$ node's label in the last Stacked-DRG layer.

$\bi K: \Label^{[N_\nodes]} = \Labels[N_\layers - 1][:]$
$\bi K_i: \Label_{i, l_\last} = \Labels[N_\layers - 1][i]$

$D$ is encoded into $R$ via *node-wise* field addition. Each $D_i \in D$ is interpreted as a field element and encoded into $R_i$ by adding $K_i$ to $D_i$. The resulting array of field elements produced via field addition is the encoding $R$ of $D$.

$\bi R: \Fq^{[N_\nodes]} = [R_0, \ldots, R_{N_\nodes - 1}]$
$\bi R_i: \Fq = D_i \as \Fqsafe \oplus K_i$

<br />

The function $\encode$ is used to encode a sector $D$ into $R$ given a an encoding key $K$ derived from $R$'s $\ReplicaID$.

$\overline{\underline{\Function \encode(D: \Safe, K: \Label^{[N_\nodes]}) \rightarrow R}}$
$\line{1}{\bi}{R: \Fq^{[N_\nodes]} = [\ ]}$
$\line{2}{\bi}{\for i \in [N_\nodes]:}$
$\line{3}{\bi}{\quad D_i: \Safe = D[i]}$
$\line{4}{\bi}{\quad K_i: \Label = K[i]}$
$\line{5}{\bi}{\quad R_i = D_i \as \Fqsafe \oplus K_i}$
$\line{6}{\bi}{\quad R\dot\push(R_i)}$
$\line{7}{\bi}{\return R}$

## Replication

Replication is the entire process by which a sector $D$ is uniquely encoded into a replica $R$. Replication encompasses Stacked-DRG labeling, encoding $D$ into $R$, and the generation of trees $\TreeC$ over $\Labels$ and $\TreeR$ over $R$.

A miner derives a unique $\ReplicaID$ for each $R$ using the commitment to the replica's sector $\CommD = \TreeD\dot\root \thin$ (where $\TreeD$ is build over the nodes of the unencoded sector $D$ associated with $R \thin$).

Given a sector $D$ and its commitment $\CommD$, replication proceeds as follows:

1. Generate the $R$'s unique $\ReplicaID$.
2. Generate $\Labels$ from $\ReplicaID$, thus deriving the key $K$ that encodes $D$ into $R$.
3. Generate $\TreeC$ over the columns of $\Labels$ via the column commitment process.
4. Encode $D$ into $R$ using the encoding key $K$.
5. Generate a $\TreeR: \OctTree$ over the replica $R$.
6. Commit to $R$ and its associated labeling $\Labels$ via the commitment $\CommCR$.

The function $\replicate$ runs the entire replication process for a sector $D$.

$\overline{\Function \replicate( \qquad\qquad\qquad\qquad\qquad\quad\bi\ }$
$\quad D: \Safe^{[N_\nodes]},$
$\quad \CommD: \Safe,$
$\quad \SectorID_D: \u{64},$
$\quad \R_\replicaid: \Byte^{[32]},$
$\underline{) \rightarrow \ReplicaID, R, \TreeC, \TreeR, \CommCR, \Labels}$
$\line{1}{\bi}{\ReplicaID = \createreplicaid(\ProverID, \SectorID, \R_\replicaid, \CommD, \PorepID)}$
$\line{2}{\bi}{\Labels = \generatelabels(\ReplicaID)}$
$\line{3}{\bi}{\TreeC = \createtreecfromlabels(\Labels)}$
$\line{4}{\bi}{K: \Label^{[N_\nodes]} = \Labels[N_\layers - 1][:]}$
$\line{5}{\bi}{R: \Fq^{[N_\nodes]} = \textsf{encode}(D, K)}$
$\line{6}{\bi}{\TreeR = \OctTree\cc\new(R)}$
$\line{7}{\bi}{\CommCR: \Fq = \Poseidon_2([\TreeC\dot\root, \TreeR\dot\root])}$
$\line{8}{\bi}{\return \ReplicaID, R, \TreeC, \TreeR, \CommCR, \Labels}$

## ReplicaID Generation

The function $\createreplicaid$ describes how a miner having the ID $\ProverID$ is able to generate a $\ReplicaID$ for a replica $R$ of sector $D$, where $D$ has a unique ID $\SectorID$ and commitment $\CommD$. The prover uses a unique random value $\R_\ReplicaID$ for each $\ReplicaID$ generated.

**Implementation:** [`storage_proofs::porep::stacked::vanilla::params::generate_replica_id()`](https://github.com/filecoin-project/rust-fil-proofs/blob/b40b34b5ef7e2a7b8c7e7ea9e574d900728dac45/storage-proofs/porep/src/stacked/vanilla/params.rs#L736)

$\overline{\Function \createreplicaid(\ }$
$\quad \ProverID: \Byte^{[32]},$
$\quad \SectorID: \u{64},$
$\quad \R_\replicaid: \Byte^{[32]},$
$\quad \CommD: \Safe,$
$\quad \PorepID: \Byte^{[32]},$
$\underline{) \rightarrow \ReplicaID \qquad\qquad\qquad\bi\ }$
$\line{1}{}{\preimage: \Byte^{[136]} =}$
$\quad\quad \ProverID$
$\quad\quad \|\ \beencode(\SectorID) \as \Byte^{[8]}$
$\quad\quad \|\ \R_\ReplicaID$
$\quad\quad \|\ \CommD$
$\quad\quad \|\ \PorepID$

$\line{2}{}{\return \Sha{254}(\preimage) \as \Fqsafe}$

## Sector Construction

A sector $D$ is constructed from Filecoin client data where the aggregating of client data of has been preprocessed/bit-padded such that two zero bits are placed between each distinct 254-bit slice of client data. This padding process results in a sector $D$ such that every 256-bit slice represents a valid 254-bit field element $\Safe \thin$.

A Merkle tree $\TreeD: \BinTree$ is constructed for sector $D$ whose leaves are the 256-bit slices $D_i: \Safe \in D \thin$.

$\bi D_i: \Safe = D[i * 32 \thin\ldotdot\thin (i + 1) * 32]$
$\bi \TreeD = \BinTree\cc\new([D_0, \ldots, D_{N_\nodes - 1}])$
$\bi \CommD: \Safe = \TreeD\dot\root$

Each $\TreeD$ is constructed over the preprocessed sector data $D$.

## PoRep Challenges

The function $\getporepchallenges$ creates the PoRep challenge set for a replica $R$'s partition-$k$ PoRep partition proof.

**Implementation:** [`storage_proofs::porep::stacked::vanilla::challenges::LayerChallenges::derive_internal()`](https://github.com/filecoin-project/rust-fil-proofs/blob/c58918b9b2f749a5db40a7952d29a6501c765e13/storage-proofs/porep/src/stacked/vanilla/challenges.rs#L40)

$\overline{\Function\ \getporepchallenges( \quad}$
$\quad \ReplicaID,$
$\quad \R_\porepchallenges: \Byte^{[32]},$
$\quad k: [N_\poreppartitions],$
$\underline{) \rightarrow \PorepChallenges_{R, k} \qquad\qquad\qquad}$
$\line{1}{\bi}{\challenges: \NodeIndex^{[N_{\porepchallenges / k}]} = [\ ]}$
$\line{2}{\bi}{\for \challengeindex \in [N_\porepchallenges]:}$
$\line{3}{\bi}{\quad \challengeindex_\porepbatch: \u{32} = k * N_{\porepchallenges / k} + \challengeindex}$
$\line{4}{\bi}{\quad \preimage: \Byte^{[68]} =}$
$\quad\quad\quad \leencode(\ReplicaID) \as \Byte^{[32]}$
$\quad\quad\quad \|\ \R_\porepchallenges$
$\quad\quad\quad \|\ \leencode(\challengeindex_\porepbatch) \as \Byte^{[4]}$
$\line{5}{\bi}{\quad \digest: \Byte^{[32]} = \Sha{256}(\preimage)}$
$\line{6}{\bi}{\quad \digestint: \u{256} = \ledecode(\digest)}$
$\line{7}{\bi}{\quad c: \NodeIndex \setminus 0 = (\digestint \MOD (N_\nodes - 1)) + 1}$
$\line{8}{\bi}{\quad \challenges\dot\push(c)}$
$\line{9}{\bi}{\return \challenges}$

## Vanilla PoRep

### Proving

A PoRep prover generates a $\PorepPartitionProof_k$ for each partition $k$ in a replica $R$'s batch of PoRep proofs. Each partition proof is generated for $N_{\porepchallenges / k}$ number of challenges, the challenge set $\PorepChallenges_{R, k}$ (each partition proof's challenge set is specific to the replica $R$ and partition $k$).

A single partition proof generated by a PoRep prover shows that:

* The prover knows a valid Merkle path for $c$ in $\TreeD$ that is consistent with the public $\CommD$.
* The prover knows valid Merkle paths for $c$ in trees $\TreeC$ and $\TreeR$ which are consistent with the committed to $\CommCR$.
* The prover knows $c$'s labeling in each Stacked-DRG layer $\Column_c = \ColumnProof_c\dot\column \thin$ by hashing $\Column_c$ into a leaf in $\TreeC$ that is consistent with $\CommCR$.
* For each layer $l$ in the Stacked-DRG, the prover knows $c$'s labeling preimage $\ParentLabels$ (taken from the columns in $\ParentColumnProofs$), such that the parent labels are consistent with $\CommCR$.
* The prover knows the key $K_c$ used to encode $D_c$ into $R_c$ (where $D_c$, $K_c$, and $R_c$ were already shown to be consistent with the commitments $\CommD$ and $\CommCR$).

$\overline{\mathbf{Function:}\ \createvanillaporepproof(\ }$
$\quad k,$
$\quad \ReplicaID,$
$\quad \TreeD,$
$\quad \TreeC,$
$\quad \TreeR,$
$\quad \Labels,$
$\quad \R_\porepchallenges: \Byte^{[32]},$
$\underline{) \rightarrow \PorepPartitionProof_{R, k} \qquad\qquad\qquad}$
$\line{1}{\bi}{\PorepPartitionProof_{R, k}: \PorepChallengeProof^{\thin[N_\porepchallenges]} = [\ ]}$
$\line{2}{\bi}{\PorepChallenges_{R, k} = \getporepchallenges(\ReplicaID, \R_\porepchallenges, k)}$

$\line{3}{\bi}{\for c \in \PorepChallenges_{R, k}:}$
$\line{4}{\bi}{\quad \TreeDProof_c = \TreeD\dot\createproof(c)}$

$\line{5}{\bi}{\quad \ColumnProof_c\ \{}$
$\quad\quad\quad \column: \Labels[:][c],$
$\quad\quad\quad \leaf,\thin \root,\thin \path \Leftarrow \TreeC\dot\createproof(c),$
$\quad\quad \}$

$\line{6}{\bi}{\quad \TreeRProof_c = \TreeR\dot\createproof(c)}$

$\line{7}{\bi}{\quad \ParentColumnProofs_{\mathbf{u}_\total}: \ColumnProof^{[d_\total]} = [\ ]}$
$\line{8}{\bi}{\quad \mathbf{u}_\total: \NodeIndex^{[d_\total]} = \getallparents(c, \PorepID)}$
$\line{9}{\bi}{\quad \for u \in \mathbf{u}_\total:}$
$\line{10}{}{\quad\quad \ColumnProof_u\ \{}$
$\quad\quad\quad\quad \column: \Labels[:][u],$
$\quad\quad\quad\quad \leaf,\thin \root,\thin \path \Leftarrow \TreeC\dot\textsf{create_proof}(u),$
$\quad\quad\quad \}$
$\line{11}{}{\quad\quad \ParentColumnProofs_{\mathbf{u}_\total}\dot\push(\ColumnProof_u)}$

$\line{12}{}{\quad \PorepChallengeProof_c\ \{}$
$\quad\quad\quad\quad \TreeDProof_c,$
$\quad\quad\quad\quad \ColumnProof_c,$
$\quad\quad\quad\quad \TreeRProof_c,$
$\quad\quad\quad\quad \ParentColumnProofs_{\mathbf{u}_\total},$
$\quad\quad \}$
$\line{13}{}{\quad \PorepPartitionProof_{R, k}\dot\push(\PorepChallengeProof_c)}$
$\line{14}{}{\return \PorepPartitionProof_{R, k}}$

### Verification

**Implementation:**
* [`storage_proofs::porep::stacked::vanilla::proof_scheme::StackedDrg::verify_all_partitions()`](https://github.com/filecoin-project/rust-fil-proofs/blob/c58918b9b2f749a5db40a7952d29a6501c765e13/storage-proofs/porep/src/stacked/vanilla/proof_scheme.rs#L82)
* [`storage_proofs::porep::stacked::vanilla::params::Proof::verify()`](https://github.com/filecoin-project/rust-fil-proofs/blob/447a8ba76da224b8b5f9b7b8dd624ba9a6a107a6/storage-proofs/porep/src/stacked/vanilla/params.rs#L191)

$\overline{\Function\ \verifyvanillaporepproof(}$
$\quad \PorepPartitionProof_{R, k} \thin,$
$\quad k: [N_\poreppartitions],$
$\quad \ReplicaID,$
$\quad \CommD,$
$\quad \CommCR,$
$\quad \R_\porepchallenges: \Byte^{[32]},$
$\underline{) \qquad\qquad\qquad\qquad\qquad\qquad\qquad\qquad}$
$\line{1}{\bi}{\PorepChallenges_{R, k} = \getporepchallenges(\ReplicaID, \R_\porepchallenges, k)}$

$\line{2}{\bi}{\for i \in [N_{\porepchallenges / k}]}:$
$\line{3}{\bi}{\quad c = \PorepChallenges_{R, k}[i]}$
$\line{4}{\bi}{\quad \TreeDProof_c, \ColumnProof_c, \TreeRProof_c, \ParentColumnProofs_{\mathbf{u}_\total}}$
$\quad\quad\quad \Leftarrow \PorepPartitionProof_{R, k}[i]$

$\line{5}{\bi}{\quad \assert(\TreeDProof_c\dot\root = \CommD)}$

$\line{6}{\bi}{\quad \CommCR^\dagger = \Poseidon_2([\ColumnProof_c\dot\root, \TreeRProof_c\dot\root])}$
$\line{7}{\bi}{\quad \assert(\CommCR^\dagger = \CommCR)}$

$\line{8}{\bi}{\quad \assert(\calculatebintreechallenge(\TreeDProof_c\dot\path) = c)}$
$\line{9}{\bi}{\quad \assert(\calculateocttreechallenge(\ColumnProof_c) = c)}$
$\line{10}{}{\quad \assert(\calculateocttreechallenge(\TreeRProof_c) = c)}$

$\line{11}{}{\quad \assert(\bintreeproofisvalid(\TreeDProof_c))}$
$\line{12}{}{\quad \assert(\octtreeproofisvalid(\ColumnProof_c))}$
$\line{13}{}{\quad \assert(\octtreeproofisvalid(\TreeRProof_c))}$

$\line{14}{}{\quad \assert(\ColumnProof_c.\leaf = \Poseidon_{11}(\ColumnProof_c.\column))}$

$\line{15}{}{\quad \mathbf{u}_\total: \NodeIndex^{[d_\total]} = \getallparents(c, \PorepID)}$
$\line{16}{}{\quad \for p \in [d_\total]:}$
$\line{17}{}{\quad\quad u = \mathbf{u}_\total[p]}$
$\line{18}{}{\quad\quad \ColumnProof_u = \ParentColumnProofs_{\mathbf{u}_\total}[p]}$
$\line{19}{}{\quad\quad \assert(\ColumnProof_u.\root = \ColumnProof_c.\root)}$
$\line{20}{}{\quad\quad \assert(\calculateocttreechallenge(\ColumnProof_u\dot\path) = u)}$
$\line{21}{}{\quad\quad \assert(\octtreeproofisvalid(\ColumnProof_u))}$
$\line{22}{}{\quad\quad \assert(\ColumnProof_u.\leaf = \Poseidon_{11}(\ColumnProof_u.\column))}$

$\line{23}{}{\quad \for l \in [N_\layers]:}$
$\line{24}{}{\quad\quad \calculatedlabel_{c, l} = \createlabel_\V(\ReplicaID, l, c, \ParentColumnProofs_{\mathbf{u}_\total})}$
$\line{25}{}{\quad\quad \assert(\calculatedlabel_{c, l} = \ColumnProof_c.\column[l])}$

$\line{26}{}{\quad D_c = \TreeDProof.\leaf}$
$\line{27}{}{\quad {R_c}^\dagger = \TreeRProof.\leaf}$
$\line{28}{}{\quad {K_c}^\dagger = \ColumnProof_c.\column[N_\layers - 1]}$
$\line{29}{}{\quad \assert({R_c}^\dagger \ominus {K_c}^\dagger = D_c)}$

#### Verifier Labeling

**Implementation:** [`storage_proofs::porep::stacked::vanilla::labeling_proof::LabelingProof::create_label()`](https://github.com/filecoin-project/rust-fil-proofs/blob/447a8ba76da224b8b5f9b7b8dd624ba9a6a107a6/storage-proofs/porep/src/stacked/vanilla/labeling_proof.rs#L27)

**Additional Notation:**
$\createlabel_\V$
Designates the function $\createlabel$ as being used by a PoRep verifier $\V$.

$c: \NodeIndex \setminus 0$
The node index of a PoRep challenge. The first node 0 is never challenged in PoRep proofs.

$d = \big[ \tern l_c = 0 ? d_\drg : d_\total \big]$ 
The number of parents that challenge $c$ has (where $c$ is in the layer $l_c$).

$\Label_{c, l}^\dagger: \Fqsafe$
The label of the challenge node $c$ in layer $l$ calculated from the unverified $\ParentColumnProofs^\dagger$.

$\Label_{u_\exp, l - 1}: \Fqsafe$
The label of challenge $c$'s expander parent $u_\exp$ in layer $l - 1$. Expander parents come from the layer prior to $c$'s layer $l$. 

$p_\drg \in [d_\drg]$
$p_\total \in [d_\total]$
The index of a parent in $c$'s parent arrays $\mathbf{u}_\drg$ and $\mathbf{u}_\total$ respectively.

$u_\drg, u_\exp: \NodeIndex$
The node index of a DRG or expander graph parent for $c$.

$\parentlabels': \Label^{[N_\parentlabels]}$
The set of parent labels repeated until its length is $N_\parentlabels$.

$\overline{\Function: \createlabel_\V( \qquad\ }$
$\quad \ReplicaID,$
$\quad l: [N_\layers],$
$\quad c: \NodeIndex \setminus 0,$
$\quad \ParentColumnProofs_{\mathbf{u}_\total}^\dagger,$
$\underline{) \rightarrow \Label_{c, l}^\dagger \qquad\qquad\qquad\qquad\quad}$
$\line{1}{\bi}{\parentlabels: {\Label_u}^{[d]} = [\ ]}$

$\line{2}{\bi}{\for p_\drg \in [d_\drg]:}$
$\line{3}{\bi}{\quad\quad \Label_{u_\drg, l} = \ParentColumnProofs_{\mathbf{u}_\total}[p_\drg]\dot\column[l]}$
$\line{4}{\bi}{\quad\quad \parentlabels\dot\push(\Label_{u_\drg, l})}$

$\line{5}{\bi}{\if l > 0:}$
$\line{6}{\bi}{\quad \for p_\exp \in [d_\drg, d_\total - 1]:}$
$\line{7}{\bi}{\quad\quad \Label_{u_\exp, l - 1} = \ParentColumnProofs_{\mathbf{u}_\total}[p_\exp]\dot\column[l - 1]}$
$\line{8}{\bi}{\quad\quad \parentlabels\dot\push(\Label_{u_\exp, l - 1})}$

$\line{9}{\bi}{\parentlabels': \Label^{[N_\parentlabels]} = \parentlabels\dot\repeattolength(N_\parentlabels)}$

$\line{10}{}{\preimage: \Byte^{[1228]} =}$
$\quad\quad \ReplicaID$
$\quad\quad \|\ \beencode(l \as \u{32}) \as \Byte^{[4]}$
$\quad\quad \|\ \beencode(c \as \u{64}) \as \Byte^{[8]}$
$\quad\quad \big{\|}_{\Label_u \hspace{1pt} \in \hspace{1pt} \parentlabels'} \thin \Label_u \as \Byte^{[32]}$

$\line{11}{}{\return \Sha{254}(\preimage) \as \Fq}$

## PoRep Circuit

**Implementation:**
* [`storage_proofs::porep::stacked::circuit::proof::StackedCircuit::synthesize()`](https://github.com/filecoin-project/rust-fil-proofs/blob/master/storage-proofs/porep/src/stacked/circuit/proof.rs#L77)
* [`storage_proofs::porep::stacked::circuit::params::Proof::synthesize()`](https://github.com/filecoin-project/rust-fil-proofs/blob/master/storage-proofs/porep/src/stacked/circuit/params.rs#L77)

**Additional Notation:**
$\PorepPartitionProof_{R, k}$
The $k^{th}$ PoRep partition proof generated for the replica $R$.

$\treedleaf_{\auxb, c}$
The circuit value for a challenge $c$'s leaf in $\TreeD$.

$\calculatedcommd_{\auxb, c}$
The circuit value calculated for $\CommD$ using challenge $c$'s $\TreeDProof_c$.

$\column_{[\auxb], u}$
The array circuit values representing a parent $u$ of challenge $c$'s label in each Stacked-DRG layer.

$\parentcolumns_{[[\auxb]], \mathbf{u}_\total}$
An array of an array of circuit values, the allocated column for each parent $u \in \mathbf{u}_\total \thin$.

$l_c$
The challenge $c$'s layer in the Stacked-DRG.

$\parentlabelsbits_{[[\auxb + \constb, \lebytes]]}$
An array where each element is a parent $u$'s label, an array of allocated and unallocated bits $_{[\auxb + \constb]}$ having $\lebytes$ bit order.

$\calculatedlabel_{\auxb, c, l}$
The label calculated for challenge $c$ in residing in layer $l$.

$\overline{\Function \createporepcircuit(}$
$\quad \PorepPartitionProof_{R, k} \thin,$
$\quad k,$
$\quad \ReplicaID,$
$\quad \CommD,$
$\quad \CommC,$
$\quad \CommR,$
$\quad \CommCR,$
$\quad \R_\porepchallenges: \Byte^{[32]},$
$\underline{) \rightarrow \RCS \qquad\qquad\qquad\qquad\qquad}$
$\line{1}{\bi}{\cs = \RCS\cc\new()}$

$\line{2}{\bi}{\replicaid_\pubb: \CircuitVal \deq \cs\dot\publicinput(\ReplicaID)}$
$\line{3}{\bi}{\replicaidbits_{[\auxb, \Le]}: \CircuitBit^{[255]}\ \deq \lebitsgadget(\cs, \replicaid_\pubb, 255)}$
$\line{4}{\bi}{\replicaidbits_{[\auxb+\constb, \lebytes]}: \CircuitBitOrConst^{[256]} =}$
$\quad\quad \lebitstolebytes(\replicaidbits_{[\auxb, \Le]})$

$\line{5}{\bi}{\commd_\pubb: \CircuitVal \deq \cs\dot\publicinput(\CommD)}$
$\line{6}{\bi}{\commcr_\pubb: \CircuitVal \deq \cs\dot\publicinput(\CommCR)}$

$\line{7}{\bi}{\commc_\auxb: \CircuitVal \deq \cs\dot\privateinput(\CommC)}$
$\line{8}{\bi}{\commr_\auxb: \CircuitVal \deq \cs\dot\privateinput(\CommR)}$

$\line{9}{\bi}{\calculatedcommcr_\auxb: \CircuitVal\ \deq}$
$\quad\quad \poseidongadget{2}(\cs, [\commc_\auxb, \thin \commr_\auxb])$
$\line{10}{}{\cs\dot\assert(\calculatedcommcr_\auxb = \commcr_\pubb)}$

$\line{11}{}{\PorepChallenges_k = \getporepchallenges(\ReplicaID, \R_\porepchallenges, k)}$

$\line{12}{}{\for i \in [N_\porepchallenges]}:$
$\line{13}{}{\quad c: \NodeIndex = \PorepChallenges[i]}$
$\line{14}{}{\quad \TreeDProof_c, \ColumnProof_c, \TreeRProof_c, \ParentColumnProofs_{\mathbf{u}_\total}}$
$\quad\quad\quad \Leftarrow \PorepPartitionProof[i]$

$\line{15}{}{\quad \challengebits_{[\auxb, \Le]}: \CircuitBit^{[64]} \deq \lebitsgadget(\cs, c, 64)}$
$\line{16}{}{\quad \packedchallenge_\pubb: \CircuitVal\ \deq}$
$\quad\quad\quad \packbitsasinputgadget(\cs, \challengebits_{[\auxb, \Le]})$

$\line{17}{}{\quad \treedleaf_{\auxb, c} \deq \cs\dot\privateinput(\TreeDProof_c\dot\leaf)}$
$\line{18}{}{\quad \calculatedcommd_{\auxb, c}\ \deq}$
$\quad\quad\quad \bintreerootgadget(\cs, \treedleaf_{\auxb, c}\thin, \TreeDProof_c\dot\path)$
$\line{19}{}{\quad \cs\dot\assert(\calculatedcommd_{\auxb, c} = \commd_\pubb)}$

$\line{20}{}{\quad \parentcolumns_{[[\auxb]], \mathbf{u}_\total}: {\CircuitVal^{[N_\layers]}}^{[d_\total]} = [\ ]}$
$\line{21}{}{\quad \for \ColumnProof_u \in \ParentColumnProofs_{\mathbf{u}_\total}:}$
$\line{22}{}{\quad\quad \column_{[\auxb], u}: \CircuitVal^{[N_\layers]}\ \deq}$
$\quad\quad\quad\quad [\thin \cs\dot\privateinput(\label_{u, l}) \mid \forall\thin \label_{u, l} \in \ColumnProof_u\dot\column \thin]$
$\line{23}{}{\quad\quad \calculatedtreecleaf_{\auxb, u}: \CircuitVal \deq \poseidongadget{11}(\cs, \column_{[\auxb], u})}$
$\line{24}{}{\quad\quad \calculatedcommc_{\auxb, u}: \CircuitVal\ \deq}$
$\quad\quad\quad\quad \octtreerootgadget(\cs,\thin \calculatedtreecleaf_{\auxb, u},\thin \ColumnProof_u\dot\path)$
$\line{25}{}{\quad\quad \cs\dot\assert(\calculatedcommc_{\auxb, c} = \commc_\auxb)}$
$\line{26}{}{\quad\quad \parentcolumns_{[[\auxb]], \mathbf{u}_\total}\dot\push(\column_{[\auxb], u})}$

$\line{27}{}{\quad \calculatedcolumn_{[\auxb], c}: \CircuitVal^{[N_\layers]} = [\ ]}$
$\line{28}{}{\quad \for l_c \in [N_\layers]:}$
$\line{29}{}{\quad\quad \layerbits_{[\auxb, \Le]}: \CircuitBit^{[32]} \deq \lebitsgadget(\cs, l_c, 32)}$

$\line{30}{}{\quad\quad \parentlabels_{[\auxb]}: \CircuitVal^{[*]} = [\ ]}$
$\line{28}{}{\quad\quad \for p_\drg \in [d_\drg]:}$
$\line{31}{}{\quad\quad\quad \parentlabel_{\auxb, u_\drg} = \parentcolumns_{[[\auxb]], \mathbf{u}_\total}[p_\drg][l_c]}$
$\line{32}{}{\quad\quad\quad \parentlabels_{[\auxb]}\dot\push(\parentlabel_{\auxb, u_\drg})}$
$\line{33}{}{\quad\quad \if l_c > 0:}$
$\line{34}{}{\quad\quad\quad \for p_\exp \in [d_\drg, d_\total - 1]:}$
$\line{35}{}{\quad\quad\quad\quad \parentlabel_{\auxb, u_\exp} = \parentcolumns_{[[\auxb]], \mathbf{u}_\total}[p_\exp][l_c - 1]}$
$\line{36}{}{\quad\quad\quad\quad \parentlabels_{[\auxb]}\dot\push(\parentlabel_{\auxb, u_\exp})}$

$\line{37}{}{\quad\quad \parentlabelsbits_{[[\auxb + \constb, \lebytes]]}: {\CircuitBitOrConst^{[256]}}^{[d_\drg\ \text{or}\ d_\exp]} = [\ ]}$
$\line{38}{}{\quad\quad \for \parentlabel_\auxb \in \parentlabels_{[\auxb]}:}$
$\line{39}{}{\quad\quad\quad \parentlabelbits_{[\auxb, \Le]}: \CircuitBit^{[255]} \deq}$
$\quad\quad\quad\quad\quad\quad \lebitsgadget(\cs, \parentlabel_\auxb, 255)$
$\line{40}{}{\quad\quad\quad \parentlabelbits_{[\auxb + \constb, \lebytes]}: \CircuitBitOrConst^{[256]} =}$
$\quad\quad\quad\quad\quad\quad \lebitstolebytes(\parentlabelbits_{[\auxb, \Le]})$
$\line{41}{}{\quad\quad\quad \parentlabelsbits_{[[\auxb + \constb, \lebytes]]}\dot\push(\parentlabelbits_{[\auxb + \constb, \lebytes]})}$
$\line{42}{}{\quad\quad \parentlabelsbits_{[[\auxb + \constb, \lebytes]]}\dot\repeat(N_\parentlabels)}$

$\line{43}{}{\quad\quad \calculatedlabel_{\auxb, c, l}: \CircuitVal \deq \createlabelgadget(}$
$\quad\quad\quad\quad \cs,$
$\quad\quad\quad\quad \replicaidbits_{[\auxb + \constb, \lebytes]} \thin,$
$\quad\quad\quad\quad \layerbits_{[\auxb, \Le]} \thin,$
$\quad\quad\quad\quad \challengebits_{[\auxb, \Le]},$
$\quad\quad\quad\quad \parentlabelsbits_{[[\auxb + \constb, \lebytes]]} \thin,$
$\quad\quad\quad)$
$\line{44}{}{\quad\quad \calculatedcolumn_{[\auxb], c}\dot\push(\calculatedlabel_{\auxb, c, l})}$

$\line{45}{}{\quad \calculatedcommr_{\auxb, c}: \CircuitVal\ \deq}$
$\quad\quad\quad \octtreerootgadget(\cs,\thin \calculatedtreerleaf_{\auxb, c} \thin, \TreeRProof_c\dot\path)$
$\line{46}{}{\quad \cs\dot\assert(\calculatedcommr_{\auxb, c} = \commr_\auxb)}$

$\line{47}{}{\quad \encodingkey_{\auxb, c}: \CircuitVal = \calculatedcolumn_{[\auxb], c}[N_\layers - 1]}$
$\line{48}{}{\quad \calculatedtreerleaf_{\auxb, c}: \CircuitVal\ \deq}$
$\quad\quad\quad \encodegadget(\cs,\thin \treedleaf_{\auxb, c} \thin, \encodingkey_{\auxb, c})$

$\line{49}{}{\quad \calculatedtreecleaf_{\auxb, c}: \CircuitVal\ \deq}$
$\quad\quad\quad \poseidongadget{11}(\cs, \calculatedcolumn_{[\auxb], c})$
$\line{50}{}{\quad \calculatedcommc_{\auxb, c}: \CircuitVal\ \deq}$
$\quad\quad\quad \octtreerootgadget(\cs,\thin \calculatedtreecleaf_{\auxb, c} \thin, \ColumnProof_c\dot\path)$
$\line{51}{}{\quad \cs\dot\assert(\calculatedcommc_{\auxb, c} = \commc_\auxb)}$

$\line{52}{}{\return \cs}$

**Code Comments:**
* **Lines 9-10:** Computes $\CommCR^\dagger$ within the circuit from the witnessed commitments and assert that $\CommCR^\dagger$ is equal to the public input $\CommCR$.
* **Lines 15-16:** Adds the packed challenge $c$ as a public input, used when calculating each challenge $c$'s column within the circuit.
* **Lines 17-19:** Verifies $c$'s $\TreeDProof_c$ by computing $\CommD_c^\dagger$ within the circuit and asserting that it is equal to the public input $\CommD$.
* **Lines 20-26:** Allocates each of $c$'s parent's $u \in \mathbf{u}_\total$ label and checks that $u$'s $\ColumnProof_u$ is consistent with the previously verified $\CommC^\dagger \mapsto \CommC \thin$.
* **Lines 27-44:** Calculates challenge $c$'s label in each Stacked-DRG layer $l$ within the circuit using each parent's allocated column.
* **Lines 45-46:** Verifies that $c$'s $\TreeRProof_c$ is consistent with the previously verified $\CommR^\dagger \mapsto \CommR$.
* **Lines 47-48:** Checks that the calculated encoding key $K_c^\dagger$ for $c$ encodes the previously verified sector and replica tree leaves $D_c^\dagger \mapsto D_c$ into $R_c^\dagger \mapsto R_c$.
* **Lines 49-51:** Verifies $c$'s $\ColumnProof_c$ against the previously verified $\CommC$.

## PoSt Challenges

The function $\getpostchallenge$ is used to derive a Merkle challenge for a Winning or Window PoSt proof.

**Implementation:** [`storage_proofs::post::fallback::vanilla::generate_leaf_challenge()`](https://github.com/filecoin-project/rust-fil-proofs/blob/8e8306c942c22571bc784f7536f1704058c45119/storage-proofs/post/src/fallback/vanilla.rs#L214)

**Additional Notation:**
$\R_{\postchallenges, \batch \thin \aww}$
A random value used to derive the challenge set for each of a PoSt prover's partition proofs in their current Winning or Window PoSt proof batch.

$\SectorID$
The ID for the sector $D$ associated with the replica $R$ for which this Merkle challemnge is being generated.

$\challengeindex_\batch$
The unique index of a Merkle challenge across all PoSt partition proofs that a PoSt prover is generating. For all partition proofs in the same PoSt batch, every Merkle challenge across all replicas will have a unique $\challengeindex_\batch \thin$.

$\overline{\Function\getpostchallenge(\qquad\qquad}$
$\quad \R_{\postchallenges, \batch \thin \aww}: \Fq,$
$\quad \SectorID: \u{64},$
$\quad \challengeindex_\batch: \u{64},$
$\underline{) \rightarrow \NodeIndex \qquad\qquad\qquad\qquad\qquad\quad}$
$\line{1}{\bi}{\preimage: \Byte^{[48]} =}$
$\quad\quad \leencode(\R_{\postchallenges, \batch \thin \aww}) \as \Byte^{[32]}$
$\quad\quad \|\ \leencode(\SectorID) \as \Byte^{[8]}$
$\quad\quad \|\ \leencode(\challengeindex_\batch) \as \Byte^{[8]}$

$\line{2}{\bi}{\digest: \Byte^{[32]} = \Sha{256}(\preimage)}$
$\line{3}{\bi}{\digestint: \u{64} = \ledecode(\digest[\ldotdot 8])}$
$\line{4}{\bi}{\return \digestint \MOD N_\nodes}$

**Code Comments:**
* **Line 4:** modding by $N_\nodes$ takes the 64-bit $\digestint$ to a 32-bit node index $\NodeIndex$.

## Vanilla PoSt

### Proving

**Implementation:**
* [`storage_proofs::post::fallback::vanilla::FallbackPoSt::prove()`](https://github.com/filecoin-project/rust-fil-proofs/blob/8e8306c942c22571bc784f7536f1704058c45119/storage-proofs/post/src/fallback/vanilla.rs#L249)

**Additional Notation:**
$\nreplicas_k$
The number of distinct replicas that the prover has for this PoSt partition proof.

$\replicaindex_k$
$\replicaindex_\batch$
The index of a challenged replica $R$ in a partition $k$'s partition proof and the index of the challenged replica across all partition proofs that a prover is generating for batch.

$\challengeindex_R$
$\challengeindex_\batch$
The index of a Merkle challenge in a challenged replica $R$ and the index of the Merkle challenge across all partition proofs that a prover is generating for batch.

$\TreeR_R, \CommC_R, \CommCR_R, \TreeRProofs_R$
The subscript $_R$ denotes each of these values as being for the replica $R$ which is distinct within the prover's PoSt batch.

$\ell_\pad$
The number of non-distinct $\PostReplicaProof{\bf \sf s}$ that are added as padding to a PoSt prover's final partition proof in a batch.

$\overline{\Function \createvanillapostproof(\ }$
$\quad k: \mathbb{N},$
$\quad \PostReplicas_{P, k \thin \aww},$
$\quad N_{\postreplicas / k \thin \aww},$
$\quad N_{\postchallenges/R \thin \aww},$
$\quad \R_{\postchallenges \thin \aww}: \Fq,$
$\underline{) \rightarrow \PostPartitionProof_{k \thin \aww} \qquad}$
$\line{1}{\bi}{\PostPartitionProof_{k \thin \aww} = [\ ]}$
$\line{2}{\bi}{\nreplicas_k = \len(\PostReplicas_k)}$

$\line{3}{\bi}{\for \replicaindex_k \in [\nreplicas_k]}$
$\line{4}{\bi}{\quad \TreeR_R, \CommC_R, \SectorID_R \Leftarrow \PostReplicas_k[\replicaindex_k]}$
$\line{5}{\bi}{\quad \replicaindex_\batch: \u{64} = k * N_{\postreplicas / k} + \replicaindex_k}$
$\line{6}{\bi}{\quad \TreeRProofs_R: {\TreeRProof}^{\thin[N_{\postchallenges / R}]} = [\ ]}$
$\line{7}{\bi}{\quad \for \challengeindex_R \in [N_{\postchallenges / R}]:}$
$\line{8}{\bi}{\quad\quad \challengeindex_\batch: \u{64} =}$
$\quad\quad\quad\quad \replicaindex_\batch * N_{\postchallenges / R} + \challengeindex_R$
$\line{9}{\bi}{\quad\quad c = \getpostchallenge(\R_\postchallenges, \SectorID, \challengeindex_\batch)}$
$\line{10}{}{\quad\quad \TreeRProof_c = \TreeR\dot\createproof(c)}$
$\line{11}{}{\quad\quad \TreeRProofs\dot\push(\TreeRProof_c)}$
$\line{12}{}{\quad \PostPartitionProof\dot\push(\PostReplicaProof \{\thin \TreeRProofs,\thin \CommC\thin \})}$

$\line{13}{}{\ell_\textsf{pad} = N_{\postreplicas / k} - \nreplicas_k}$
$\line{14}{}{\for i \in [\ell_\textsf{pad}]}$
$\line{15}{}{\quad \PostPartitionProof\dot\push(\PostPartitionProof[\nreplicas_k - 1])}$

$\line{16}{}{\return \PartitionProof_{k \thin \aww}}$

**Code Comments:**
* **Lines 13-15:** If the prover does not have enough replicas to fill an entire PoSt partition proof, pad the partition proof with copies of the last distinct replica's $\PostReplicaProof_R \thin$.

### Verification

**Implementation:** [`storage_proofs::post::fallback::vanilla::FallbackPoSt::verify_all_partitions()`](https://github.com/filecoin-project/rust-fil-proofs/blob/8e8306c942c22571bc784f7536f1704058c45119/storage-proofs/post/src/fallback/vanilla.rs#L357)

**Additional Notation:**
$k: N_{\postpartitions, \P \thin \aww}$
The number of partitions in a Winning or Window PoSt batch is dependent on the length of the PoSt prover $\P$'s replica set.

$\nreplicas_k$
The number of distinct replicas that the prover has for this PoSt partition proof.

$\replicaindex_k$
$\replicaindex_\batch$
The index of a challenged replica $R$ in a partition $k$'s partition proofs in a PoSt prover's batch.

$\challengeindex_R$
$\challengeindex_\batch$
The index of a Merkle challenge in a challenged replica $R$ and the index of the Merkle challenge across all partition proofs in a PoSt prover's batch.

$\overline{\Function \verifyvanillapostproof( \bi}$
$\quad \PostPartitionProof_{k \thin \aww},$
$\quad k: [N_{\postpartitions, \P \thin \aww}],$
$\quad \PostReplicas_{\V, k \thin \aww},$
$\quad N_{\postreplicas / k \thin \aww},$
$\quad N_{\postchallenges / R \thin \aww},$
$\quad \R_\postchallenges: \Fq,$
$\underline{) \qquad\qquad\qquad\qquad\qquad\qquad\qquad\quad\bi}$
$\line{1}{\bi}{\nreplicas_k = \len(\PostReplicas_{\V, k})}$

$\line{2}{\bi}{\for \replicaindex_k \in [\nreplicas_k]:}$
$\line{3}{\bi}{\quad \replicaindex_\batch = k * N_{\postreplicas / k} + \replicaindex_k}$

$\line{4}{\bi}{\quad \SectorID, \CommCR \Leftarrow \PostReplicas_{\V, k}[\replicaindex_k]}$
$\line{5}{\bi}{\quad \CommC^\dagger, \TreeRProofs^\dagger \Leftarrow \PostPartitionProof[\replicaindex_k]}$
$\line{6}{\bi}{\quad \CommR^\dagger = \TreeRProofs^\dagger[0]\dot\root}$

$\line{7}{\bi}{\quad \CommCR^\dagger = \Poseidon{2}([\CommC^\dagger, \CommR^\dagger])}$
$\line{8}{\bi}{\quad \assert(\CommCR^\dagger = \CommCR)}$

$\line{9}{\bi}{\quad \for \challengeindex_R \in [N_{\postchallenges / R}]:}$
$\line{10}{}{\quad\quad \challengeindex_\batch: \u{64} =}$
$\quad\quad\quad\quad \replicaindex_\batch * N_{\postreplicas / k} + \challengeindex_R$
$\line{11}{}{\quad\quad c = \getpostchallenge(\R_\postchallenges, \SectorID, \challengeindex_\batch)}$

$\line{12}{}{\quad\quad \TreeRProof^\dagger = \TreeRProofs^\dagger[\challengeindex_R]}$
$\line{13}{}{\quad\quad \assert(\TreeRProof^\dagger\dot\root = \CommR)}$
$\line{14}{}{\quad\quad \assert(\calculateocttreechallenge(\TreeRProof^\dagger\dot\path) = c)}$
$\line{15}{}{\quad\quad \assert(\octtreeproofisvalid(\TreeRProof^\dagger))}$

**Code Comments:**
* **Line 13:** The dagger is removed from $\CommR^\dagger$ (producing $\CommR$) because $\CommR^\dagger$ was verified to be consistent with the committed to $\CommCR$ (Line 8).

## PoSt Circuit

The function $\createpostcircuit$ is used to instantiate a Winning or Window PoSt circuit.

**Addional Notation:**
$\PostPartitionProof_{k \thin \aww}$
The partition-$k$ proof in a PoSt prover's Winning or Window PoSt batch. $\PostPartitionProof_k$ Contains any padded $\PostReplicaProof{\bf \sf s}$.

$\TreeR_R, \CommC_R, \CommCR_R$
Each $\PostReplica_R \in \PostReplicas_{\P \thin \aww}$ represents a unique replica $R$ in the batch denoted by the subscript $_R \thin$.

$\TreeRProofs_R$
Each $\TreeRProofs$ is for a distinct replica $R$, denoted by the subscript $_R \thin$, in a PoSt batch.

$\overline{\Function \createpostcircuit( \quad\qquad}$
$\quad \PostPartitionProof_{k \thin \aww},$
$\quad \PostReplicas_{\P, k \thin \aww},$
$\quad N_{\postreplicas / k \thin \aww},$
$\underline{) \rightarrow \RCS \qquad\qquad\qquad\qquad\qquad\qquad\bi}$
$\line{1}{\bi}{\cs = \RCS\cc\new()}$
$\line{2}{\bi}{\nreplicas_k = \len(\PostReplicas_{\P, k})}$

$\line{3}{\bi}{\for \replicaindex_k \in [\nreplicas_k]:}$
$\line{4}{\bi}{\quad \TreeR_R, \CommC_R, \CommCR_R \Leftarrow \PostReplicas_{\P, k}[\replicaindex_k]}$
$\line{5}{\bi}{\quad \TreeRProofs_R \Leftarrow \PostPartitionProof_k[\replicaindex_k]}$

$\line{6}{\bi}{\quad \commcr_\pubb: \CircuitVal \deq \cs\dot\publicinput(\CommCR)}$
$\line{7}{\bi}{\quad \commc_\auxb: \CircuitVal \deq \cs\dot\privateinput(\CommC)}$
$\line{8}{\bi}{\quad \commr_\auxb: \CircuitVal \deq \cs\dot\privateinput(\CommR)}$
$\line{9}{\bi}{\quad \calculatedcommcr_\auxb: \CircuitVal\ \deq}$
$\quad\quad\quad \poseidongadget{2}(\cs, [\commc_\auxb, \thin \commr_\auxb])$
$\line{10}{}{\quad \cs\dot\assert(\calculatedcommcr_\auxb = \commcr_\pubb)}$

$\line{11}{}{\quad \for \TreeRProof_c \in \TreeRProofs:}$
$\line{12}{}{\quad\quad \treerleaf_{\auxb, c}: \CircuitVal \deq \cs\dot\privateinput(\TreeRProof_c\dot\leaf)}$
$\line{13}{}{\quad\quad \calculatedcommr_{\auxb, c}: \CircuitVal\ \deq}$
$\quad\quad\quad\quad \octtreerootgadget(\cs,\thin \treerleaf_{\auxb, c} \thin, \TreeRProof_c\dot\path)$
$\line{14}{}{\quad\quad \cs\dot\assert(\calculatedcommr_{\auxb, c} = \commr_\auxb)}$

$\line{15}{}{\return \cs}$

## Gadgets

### Hash Functions

We make use of the following hash function gadgets, however their implementation is beyond the scope of this document.

$\textsf{sha256_gadget}(\cs: \RCS,\thin \preimage: \CircuitBitOrConst^{[*]}) \rightarrow \CircuitBit^{[256]}$
$\shagadget{254}{2}(\cs: \RCS,\thin \inputs: \CircuitVal^{[2]}) \rightarrow \CircuitBit^{[254]}$
$\poseidongadget{2}(\cs: \RCS,\thin \inputs: \CircuitVal^{[2]}) \rightarrow \CircuitVal$
$\poseidongadget{8}(\cs: \RCS,\thin \inputs: \CircuitVal^{[8]}) \rightarrow \CircuitVal$
$\poseidongadget{11}(\cs: \RCS,\thin \inputs: \CircuitVal^{[11]}) \rightarrow \CircuitVal$

### BinTree Root Gadget

The function $\bintreerootgadget$ calculates and returns a $\BinTree$ Merkle root from an allocated leaf $\leaf_\auxb$ and an unallocated Merkle $\path$. Both the leaf and path are from a Merkle challenge $c$'s proof $\BinTreeProof_c$, where $\path = \BinTreeProof_c\dot\path \thin$.

The gadget adds one public input to the constraint system for the packed Merkle proof path bits $\pathbits_\auxle$ which are the binary representation of the $c$'s DRG node-index $\llcorner c \lrcorner_{\lower{2pt}{2, \Le}} \equiv \pathbits_\auxle \thin$).

$\overline{\Function \bintreerootgadget(\qquad\quad}$
$\quad \cs: \RCS,$
$\quad \leaf_\auxb: \CircuitVal,$
$\quad \path: \BinPathElement^{[\BinTreeDepth]},$
$\underline{) \rightarrow \CircuitVal \qquad\qquad\qquad\qquad\qquad\quad}$
$\line{1}{\bi}{\curr_\auxb: \CircuitVal = \leaf_\auxb}$
$\line{2}{\bi}{\pathbits_{[\auxb, \Le]}: \CircuitBit^{[\BinTreeDepth]} = [\ ]}$

$\line{3}{\bi}{\for \sibling, \missing \in \path:}$
$\line{4}{\bi}{\quad \missingbit_\auxb: \CircuitBit \deq \cs\dot\privateinput(\missing)}$
$\line{5}{\bi}{\quad \sibling_\auxb: \CircuitVal \deq \cs\dot\privateinput(\sibling)}$
$\line{6}{\bi}{\quad \inputs_{[\auxb]}: \CircuitVal^{[2]} \deq}$
$\quad\quad\quad \insertgadget{2}(\cs, [\sibling_\auxb],\thin \curr_\auxb,\thin \missingbit_\auxb)$
$\line{7}{\bi}{\quad \curr_\auxb: \CircuitVal \deq \shagadget{254}{2}(\cs, \inputs_{[\auxb]})}$
$\line{8}{\bi}{\quad \pathbits_{[\auxb, \Le]}\dot\push(\missingbit_\auxb)}$

$\line{9}{\bi}{\packedchallenge_\pubb: \CircuitVal\ \deq}$
$\quad\quad \packbitsasinputgadget(\cs, \pathbits_{[\auxb, \Le]})$

$\line{10}{}{\return \curr_\auxb}$

**Code Comments:**
* **Line 9:** A public input is added to $\cs$ for the Merkle challenge $c$ corresponding to the Merkle path which was used to calculate the returned root.
* **Line 10:** The final value for $\curr_\auxb$ is the Merkle root calculated from $\leaf_\auxb$ and $\path$.

### OctTree Root Gadget

The function $\octtreerootgadget$ calculates and returns an $\OctTree$ Merkle root from an allocated leaf $\leaf_\auxb$ and an unallocated Merkle $\path$. Both the leaf and path are from a Merkle challenge $c$'s proof $\OctTreeProof_c$, where $\path = \OctTreeProof_c\dot\path \thin$.

The gadget adds one public input to the constraint system for the packed Merkle proof path bits $\pathbits_\auxle$ which are the binary representation of the $c$'s DRG node-index $\llcorner c \lrcorner_{\lower{2pt}{2, \Le}} \equiv \pathbits_\auxle \thin$).

Note that the constant $3 = \log_2(8)$, the number of bits required to represent an index in the 8-element Merkle hash $\inputs$ array, is used at various times in the following algorithm.

$\overline{\Function \octtreerootgadget( \qquad\quad}$
$\quad \cs: \RCS,$
$\quad \leaf_\auxb: \CircuitVal,$
$\quad \path: \OctPathElement^{[\OctTreeDepth]},$
$\underline{) \rightarrow \CircuitVal \qquad\qquad\qquad\qquad\qquad\quad}$
$\line{1}{\bi}{\curr_\auxb: \CircuitVal = \leaf_\auxb}$
$\line{2}{\bi}{\pathbits_\auxle: \CircuitBit^{[3 * \OctTreeDepth]} = [\ ]}$

$\line{3}{\bi}{\for \siblings, \missing \in \path:}$
$\line{4}{\bi}{\quad \missingbits_\auxle: \CircuitBit^{[3]} = [\ ]}$
$\line{5}{\bi}{\quad \for i \in [3]:}$
$\line{6}{\bi}{\quad\quad \bit: \Bit = (\missing \gg i) \AND 1}$
$\line{7}{\bi}{\quad\quad \bit_\auxb \deq \cs\dot\privateinput(\bit)}$
$\line{8}{\bi}{\quad\quad \missingbits_\auxle\dot\push(\bit_\auxb)}$

$\line{9}{\bi}{\quad \siblings_{[\auxb]}: \CircuitVal^{[7]} = [\ ]}$
$\line{10}{}{\quad \for \sibling \in \siblings:}$
$\line{11}{}{\quad\quad \sibling_\auxb: \CircuitVal \deq \cs\dot\privateinput(\sibling)}$
$\line{12}{}{\quad\quad \siblings_{[\auxb]}\dot\push(\sibling_\auxb)}$

$\line{13}{}{\quad \inputs_{[\auxb]}: \CircuitVal^{[8]}\thin \deq}$
$\quad\quad\quad \insertgadget{8}(\cs, \siblings_{[\auxb]},\thin \curr_\auxb,\thin \missingbits_\auxle)$
$\line{14}{}{\quad \curr_\auxb: \CircuitVal \deq \poseidongadget{8}(\cs, \inputs_{[\auxb]})}$
$\line{15}{}{\quad \pathbits_\auxle\dot\extend(\missingbits_\auxle)}$

$\line{16}{}{\packedchallenge_\pubb: \CircuitVal\ \deq}$
$\quad\quad \packbitsasinputgadget(\cs, \pathbits_{[\auxb, \Le]})$

$\line{17}{}{\return \curr_\auxb}$

**Code Comments:**
* **Line 1:** Not a reallocation of $\leaf_\auxb$ within $\cs$, but is an in-memory copy.
* **Lines 4-8:** Witnesses the 3-bit missing index for each path element. The first iteration $i = 0$ corresponds to the least significant bit in $\missing$.
* **Lines 9-12:** Witnesses each path element's 7 Merkle hash inputs (the exlucded 8-th Merkle hash input is the calculated hash input $\curr_\auxb$ for this tree depth).
* **Line 13:** Creates the Merkle hash inputs array by inserting $\curr$ into $\siblings$ at index $\missing$.
* **Line 14:** Hashes the 8 Merkle hash inputs.
* **Line 16:** Adds the challenge $c$ as a public input.
* **Line 17:** Returns the calculated root.

### Encoding Gadget

The function $\encodegadget$ runs the $\encode$ function within a circuit. Used to encode $\unencoded_{\auxb, v}$ (node $v$'s sector data $D_v$) into $\encoded_{\auxb, v}$ (the replica node $R_v$) given an allocated encoding key $\key_{\auxb, v}$ ($K_v$).

**Implementation:** [`storage_proofs::core::gadgets::encode::encode()`](https://github.com/filecoin-project/rust-fil-proofs/blob/3cc8b5acb742d1469e10949fadc389806fa19c8c/storage-proofs/core/src/gadgets/encode.rs#L7)

$\overline{\Function \encodegadget(\qquad}$
$\quad \cs: \RCS,$
$\quad \unencoded_{\auxb, v}: \CircuitVal,$
$\quad \key_{\auxb, v}: \CircuitVal,$
$\underline{) \rightarrow \CircuitVal \qquad\qquad\qquad\qquad}$
$\line{1}{\bi}{R_v: \Fq = \unencoded_{\auxb, v}\dot\value \oplus \key_{\auxb, v}\dot\value}$
$\line{2}{\bi}{\encoded_{\auxb, v}: \CircuitVal \deq \cs\dot\privateinput(R_v)}$
$\line{3}{\bi}{\lc_A: \LinearCombination \equiv \unencoded_{\auxb, v} + \key_{\auxb, v}}$
$\line{4}{\bi}{\lc_B: \LinearCombination \equiv \cs\dot\one_\pubb}$
$\line{5}{\bi}{\lc_C: \LinearCombination \equiv \encoded_{\auxb, v}}$
$\line{6}{\bi}{\cs\dot\assert(\lc_A * \lc_B = \lc_C)}$
$\line{7}{\bi}{\return \encoded_{\auxb, v}}$

### Labeling Gadget

The function $\createlabelgadget$ is used to label a node $\node$ in the Stacked-DRG layer $\layerindex$ given the node's expanded parent labels $\parentlabels$.

**Implementation:** [`storage_proofs::porep::stacked::circuit::create_label::create_label_circuit()`](https://github.com/filecoin-project/rust-fil-proofs/blob/3cc8b5acb742d1469e10949fadc389806fa19c8c/storage-proofs/porep/src/stacked/circuit/create_label.rs#L10)

**Additional Notation:**
$\replicaid_{[\auxb + \constb, \lebytes]}$
The allocated bits (and constant zero bit(s)) representing a $\ReplicaID$.

$\layerindex_{[\auxb, \Le]}$
The allocated bits representing a layer $l \in [N_\layers]$ as an unsigned 32-bit integer.

$\node_{[\auxb, \Le]}$
A node index $v \in [N_\nodes]$ allocated as 64 bits.

$\parentlabels_{[[\auxb + \constb, \lebytes]]}$
An array containing $N_\parentlabels$ allocated bit arrays, where each bit array is the label of one of $\node$'s parents.

$\label$
Is the calculated label for $\node$.

$\overline{\Function \createlabelgadget( \qquad\qquad\qquad\qquad\qquad\qquad\qquad\qquad\ }$
$\quad \cs: \RCS,$
$\quad \replicaid_{[\auxb + \constb, \lebytes]}: \CircuitBitOrConst^{[256]} \thin,$
$\quad \layerindex_{[\auxb, \Le]}: \CircuitBit^{[32]} \thin,$
$\quad \node_{[\auxb, \Le]}: \CircuitBit^{[64]} \thin,$
$\quad \parentlabels_{[[\auxb + \constb, \lebytes]]}: {\CircuitBitOrConst^{[256]}}^{[N_\parentlabels]} \thin,$
$\underline{) \rightarrow \CircuitVal \qquad\qquad\qquad\qquad\qquad\qquad\qquad\qquad\qquad\qquad\qquad\qquad}$
$\line{1}{\bi}{\layerindex_{[\auxb, \be]}: \CircuitBit^{[32]} = \reverse(\layerindex_{[\auxb, \Le]})}$
$\line{2}{\bi}{\nodeindex_{[\auxb, \be]}: \CircuitBit^{[64]} = \reverse(\node_{[\auxb, \Le]})}$
$\line{3}{\bi}{\preimage_{[\auxb + \constb]}: \CircuitBitOrConst^{[9984]} =}$
$\quad\quad \replicaid_{[\auxb + \constb, \lebytes]}$
$\quad\quad \|\ \layerindex_{[\auxb, \be]}$
$\quad\quad \|\ \nodeindex_{[\auxb, \be]}$
$\quad\quad \|\ 0^{[160]}$
$\quad\quad \big{\|}_{\parentlabel \hspace{1pt} \in \hspace{1pt} \parentlabels} \thin \parentlabel_{[\auxb + \constb, \lebytes]} \vphantom{|^|}$

$\line{4}{\bi}{\digestbits_{[\auxb, \lebytes]}: \CircuitBit^{[256]} \deq \textsf{sha256_gadget}(\cs, \preimage_{[\auxb + \constb]})}$
$\line{5}{\bi}{\digestbits_{[\auxb, \Le]}: \CircuitBit^{[256]} = \lebytestolebits(\digestbits_{[\auxb, \lebytes]})}$
$\line{6}{\bi}{\digestbits_{[\auxb, \Le], \safe}: \CircuitBit^{[254]} = \digestbits_{[\auxb, \Le]}[0 \thin\ldotdot\thin 254]}$
$\line{7}{\bi}{\label = \digestbits_{[\auxb, \Le], \safe} \thin\as\thin \Fqsafe}$
$\line{8}{\bi}{\label_\auxb: \CircuitVal \deq \cs\dot\privateinput(\label)}$

$\line{9}{\bi}{\lc: \LinearCombination \equiv \sum_{i \in [254]}{2^i * \digestbits_{[\auxb, \Le], \safe}[i]}}$
$\line{10}{}{\cs\dot\assert(\lc = \label_\auxb)}$

$\line{11}{}{\return \label_\auxb}$

**Code Comments:**
* **Line 3:** The constant $9984 = (2 + N_\parentlabels) * \ell_\block^\bit = (2 + 37) * 256 \thin$. The constant $160 = \ell_\block^\bit - \len(\layerindex) - \len(\nodeindex) = 256 - 32 - 64 \thin$. 
* **Lines 4-5:** The constant $256 = \ell_\block^\bit \thin$.
* **Lines 5-6:** These are not reallocations.
* **Lines 6-7:** The labeling function is $\Sha{254}$ not $\Sha{256}$.
* **Lines 6,9:** The constant $254 = \ell_{\Fq, \safe}^\bit \thin$.

### Little-Endian Bits Gadget

The function $\lebitsgadget$ receives a value $\value$ allocated within a constraint system $\cs$ and reallocates it as its $n$-bit little-endian binary representation.

Note that the number of bits returned must be at least the number of bits required to represent $\value$: $0 < \lceil \log_2(\value\dot\int) \rceil \leq n \thin$.

**Implementation:** [`bellman::gadgets::num::AllocatedNum::to_bits_le()`](https://github.com/filecoin-project/bellman/blob/e0ac6b879eac87832ba6ef2e37320e877d1d96b6/src/gadgets/num.rs#L193)

$\overline{\Function \lebitsgadget(}$
$\quad \cs: \RCS,$
$\quad \value_{\langle \auxb | \pubb \rangle}: \CircuitVal,$
$\quad n: \mathbb{Z}^+,$
$\underline{) \rightarrow \CircuitBit^{[n]} \qquad\quad}$
$\line{1}{\bi}{\assert(n \geq \lceil \log_2(\value\dot\int) \rceil)}$

$\line{2}{\bi}{\bits_\Le: \Bit^{[n]} = \lower{2pt}{\llcorner} \value_{\langle \auxb | \pubb \rangle}\dot\int \lower{2pt}{\lrcorner_{2, \Le}}}$
$\line{3}{\bi}{\bits_{[\auxb, \Le]}: \CircuitBit^{[n]} = [\ ]}$
$\line{4}{\bi}{\for \bit \in \bits_\Le:}$
$\line{5}{\bi}{\quad \bit_\auxb: \CircuitBit \overset{\diamond}{=} \cs\dot\privateinput(\bit)}$
$\line{6}{\bi}{\quad \bits_{[\auxb, \Le]}\dot\push(\bit_\auxb)}$

$\line{7}{\bi}{\lc: \LinearCombination \equiv \sum_{i \in [n]}{2^i * \bits_{[\auxb, \Le]}[i]}}$
$\line{8}{\bi}{\cs\dot\assert(\value_{\langle \auxb | \pubb \rangle} = \lc)}$

$\line{9}{\bi}{\return \bits_{[\auxb, \Le]}}$

**Code Comments:**
* **Line 2:** This will pad $n - \lceil \log_2(\value\dot\int) \rceil$ zero bits onto the most significant end of $\lower{2pt}{\llcorner} \int \lower{2pt}{\lrcorner_{2, \Le}} \thin$.

### Pack Bits as Input Gadget

The function $\packbitsasinputgadget$ receives an array of $n$ allocated little-endian bits $\bits_{[\auxb, \Le]}$, where $0 < n \leq \ell_\Fqsafe^\bit \thin$, and creates the field element $\packed$ whose little-endian binary representation is that of $\bits$. The gadget adds one public input $\packed_\pubb$ to the constraint system for the created field element.

$\overline{\Function \packbitsasinputgadget(}$
$\quad \cs: \RCS,$
$\quad \bits_{[\auxb, \Le]}: \CircuitBit^{[n]},$
$\underline{) \rightarrow \CircuitVal \qquad\qquad\qquad\qquad\qquad\quad}$
$\line{1}{\bi}{\assert(0 < n \leq \ell_\Fqsafe^\bit)}$
$\line{2}{\bi}{\packed: \Fq = \bits_{[\auxb, \Le]} \as \Fq}$
$\line{3}{\bi}{\packed_\pubb \overset{\diamond}{=} \cs\dot\publicinput(\packed)}$
$\line{4}{\bi}{\lc: \LinearCombination \equiv \sum_{i \in [n]}{2^i * \bits_{[\auxb, \Le]}[i]}}$
$\line{5}{\bi}{\cs\dot\assert(\lc = \packed_\pub)}$
$\line{6}{\bi}{\return \packed_\pubb}$

### Pick Gadget

The $\pickgadget$ is used to choose one of two allocated values, $\x$ and $\y$, based upon the value of a binary condition $\bit$.

If $\bit$ is set, the gadget will reallocate and return $\x$, otherwise if $\bit$ is not set, the gadget will reallocate and return $\y$.

The $\pickgadget$, when given two allocated values $\x, \y \in \Fq$ and an allocated boolean constrained value $\bit \in \Bit$, outputs the allocated value $\pick \in \{ \x, \y \}$ and adds the $\RCS$ quadratic constraint:

$\bi (\y - \x) * (\bit) = (\y - \pick)$

This table shows that for $\bit \in \Bit$ and $\x, \y \in \Fq$ that the constraint is satisfied for the outputted values of $\pick$.

| $\bit$ | $\pick$ | $(\y - \x) * (\bit) = (\y - \pick)$ |
|--------|---------|------------------------------------|
| $1$    | $\x$     | $(\y-\x) * (1) = (\y-\x)$ |
| $0$    | $\y$     | $(\y-\x) * (0) = (\y-\y)$ |

<br />

$\overline{\Function \pickgadget( \qquad\quad\bi}$
$\quad \cs: \RCS,$
$\quad \bit_\aap: \CircuitBit,$
$\quad \x_\aap: \CircuitVal,$
$\quad \y_\aap: \CircuitVal,$
$\underline{) \rightarrow \CircuitVal \qquad\qquad\qquad\qquad}$
$\line{1}{\bi}{\pick_\auxb: \CircuitVal \deq \if \bit_\aap\dot\int = 1:}$
$\quad\quad \cs\dot\privateinput(\x_\aap)$
$\quad\else:$
$\quad\quad \cs\dot\privateinput(\y_\aap)$

$\line{2}{\bi}{\lc_A: \LinearCombination \equiv \y_\aap - \x_\aap}$
$\line{3}{\bi}{\lc_B: \LinearCombination \equiv \bit_\aap}$
$\line{4}{\bi}{\lc_C: \LinearCombination \equiv \y_\aap - \pick_\auxb}$
$\line{5}{\bi}{\cs\dot\assert(\lc_A * \lc_B = \lc_C)}$

$\line{6}{\bi}{\return \pick_\auxb}$

### Insert-2 Gadget

The $\insertgadget{2}$ inserts $\value$ into an array $\arr$ at index $\index$ and returns the inserted array of reallocated elements.

The gadget receives an array containing one allocated element $\arr[0]$ and a second allocated value $\value$ and returns the two element array containing the reallocations of the two values where the index of the reallocated $\value$ is at the index $\index$ argument in the returned 2-element array.

$\overline{\Function \insertgadget{2}(\qquad\qquad\bi}$
$\quad \cs: \RCS,$
$\quad \arr_\aap: \CircuitVal^{[1]},$
$\quad \value_\aap: \CircuitVal,$
$\quad \index_\auxb: \CircuitBitOrConst,$
$\underline{) \rightarrow \CircuitVal^{[2]} \qquad\qquad\qquad\qquad\qquad}$
$\line{1}{\bi}{\el_{\auxb, 0}: \CircuitVal \deq \pickgadget(\cs,\thin \index_\auxb,\thin \arr_\aap[0] \thin,\thin \value_\aap \thin)}$
$\line{2}{\bi}{\el_{\auxb, 1}: \CircuitVal \deq \pickgadget(\cs, \index_\auxb,\thin \value_\aap \thin,\thin \arr_\aap[0] \thin)}$
$\line{3}{\bi}{\return [\el_{\auxb, 0}, \el_{\auxb, 1}]}$

### Insert-8 Gadget

The function $\insertgadget{8}$ inserts a value $\value$ into an array of 7 elements $\arr$ at index in the 8 element array given by $\indexbits$. The values returned in the 8-element array are reallocations of $\arr$ and $\value$.

**Implementation:** [`storage_proofs::core::gadgets::insertion::insert_8()`](https://github.com/filecoin-project/rust-fil-proofs/blob/b9126bf56cfd2c73ce4ce1f6cb46fe001450f326/storage-proofs/core/src/gadgets/insertion.rs#L170)

Note that the length of the $\indexbits$ argument is $3 = \log_2(8)$ bits, which is the number of bits required to represent an index in an array of 8 elements.

**Additional Notation:**
$\arr'$
The inserted array containing 8 reallocated values, the elements of the uninserted array $\arr$ and the insertion value $\value$.

$\nor_{\auxb \thin (b_0, b_1)}$
Set to true if neither $\indexbits[0]$ nor $\indexbits[1]$ are $1$.

$\and_{\auxb \thin (b_0, b_1)}$
Set to true if both $\indexbits[0]$ and $\indexbits[1]$ are $1$.

$\pick_{\auxb, i(b_0)}$
$\pick_{\auxb, i(b_0, b_1)}$
$\pick_{\auxb, i(b_0, b_1, b_2)}$
The pick for the $i^{th}$ element of the inserted array based upon the value of the first bit (least-significant), first and second bits, and the first, second and third bits respectively. 

$b_i \equiv \indexbits_{[\aap, \Le]}[i]$
$\pick_{i(b_0)} \equiv \pick_{\auxb, i(b_0)}$
$\nor_{(b_0, b_1)} \equiv \nor_{\auxb \thin (b_0, b_1)}$
$\and_{(b_0, b_1)} \equiv \and_{\auxb \thin (b_0, b_1)}$
$\arr[i] \equiv \arr_{[\aap]}[i]$
$\arr'[i] \equiv \arr'_{[\auxb]}[i]$
For ease of notation the subscripts $_\auxb$ and $_\aap$ are left off everywhere except in the function signature and when allocating of a value within the circuit.

$\overline{\Function \insertgadget{8}( \qquad\qquad\qquad\qquad\qquad\quad}$
$\quad \cs: \RCS,$
$\quad \arr_{[\aap]}: \CircuitVal^{[7]},$
$\quad \value_\aap: \CircuitVal,$
$\quad \indexbits_{[\aap, \Le]}: \CircuitBitOrConst^{[3]},$
$\underline{) \rightarrow \CircuitVal^{[8]} \qquad\qquad\qquad\qquad\qquad\qquad\qquad\qquad\bi}$
$\line{1}{\bi}{\nor_{(b_0, b_1)}: \CircuitBit \deq \norgadget(\cs,\thin b_0 \thin,\thin b_1)}$
$\line{2}{\bi}{\and_{\auxb \thin (b_0, b_1)}: \CircuitBit \deq \andgadget(\cs,\thin b_0 \thin,\thin b_1)}$

$\line{3}{\bi}{\arr'_{[\auxb]}: \CircuitVal^{[8]} = [\ ]}$

$\line{4}{\bi}{\pick_{\auxb, 0(b_0, b_1)}: \CircuitVal \deq \pickgadget(\cs,\thin \nor_{(b_0, b_1)},\thin \value,\thin \arr[0])}$
$\line{5}{\bi}{\pick_{\auxb, 0(b_0, b_1, b_3)}: \CircuitVal \deq \pickgadget(\cs,\thin b_2,\thin \arr[0],\thin \pick_{0(b_0, b_1)})}$
$\line{6}{\bi}{\arr'[0] = \pick_{0(b_0, b_1, b_3)}}$

$\line{7}{\bi}{\pick_{\auxb, 1(b_0)}: \CircuitVal \deq \pickgadget(\cs,\thin b_0,\thin \value,\thin \arr[0])}$
$\line{8}{\bi}{\pick_{\auxb, 1(b_0, b_1)}: \CircuitVal \deq \pickgadget(\cs,\thin b_1,\thin \arr[1],\thin \pick_{1(b_0)})}$
$\line{9}{\bi}{\pick_{\auxb, 1(b_0, b_1, b_3)}: \CircuitVal \deq \pickgadget(\cs,\thin b_2,\thin \arr[1],\thin \pick_{1(b_0, b_1)})}$
$\line{10}{}{\arr'[1] = \pick_{1(b_0, b_1, b_3)}}$

$\line{11}{}{\pick_{\auxb, 2(b_0)}: \CircuitVal \deq \pickgadget(\cs,\thin b_0,\thin \arr[2],\thin \value)}$
$\line{12}{}{\pick_{\auxb, 2(b_0, b_1)}: \CircuitVal \deq \pickgadget(\cs,\thin b_1,\thin \pick_{2(b_0)},\thin \arr[1])}$
$\line{13}{}{\pick_{\auxb, 2(b_0, b_1, b_2)}: \CircuitVal \deq \pickgadget(\cs,\thin b_2,\thin \arr[2],\thin \pick_{2(b_0, b_1)})}$
$\line{14}{}{\arr'[2] = \pick_{2(b_0, b_1, b_3)}}$

$\line{15}{}{\pick_{\auxb, 3(b_0, b_1)}: \CircuitVal \deq \pickgadget(\cs, \thin \and_{(b_0, b_1)}, \thin \value, \thin \arr[2])}$
$\line{16}{}{\pick_{\auxb, 3(b_0, b_1, b_2)}: \CircuitVal \deq \pickgadget(\cs, \thin b_2, \thin \arr[3], \thin \pick_{3(b_0, b_1)})}$
$\line{17}{}{\arr'[3] = \pick_{3(b_0, b_1, b_3)}}$

$\line{18}{}{\pick_{\auxb, 4(b_0, b_1)}: \CircuitVal \deq \pickgadget(\cs, \thin \nor_{(b_0, b_1)}, \thin \value, \thin \arr[4])}$
$\line{19}{}{\pick_{\auxb, 4(b_0, b_1, b_2)}: \CircuitVal \deq \pickgadget(\cs, \thin b_2, \thin \pick_{4(b_0, b_1)}, \thin \arr[3])}$
$\line{20}{}{\arr'[4] = \pick_{4(b_0, b_1, b_3)}}$

$\line{21}{}{\pick_{\auxb, 5(b_0)}: \CircuitVal \deq \pickgadget(\cs, \thin b_0, \thin \value, \thin \arr[4])}$
$\line{22}{}{\pick_{\auxb, 5(b_0, b_1)}: \CircuitVal \deq \pickgadget(\cs, \thin b_1, \thin \arr[5], \thin \pick_{5(b_0)})}$
$\line{23}{}{\pick_{\auxb, 5(b_0, b_1, b_2)}: \CircuitVal \deq \pickgadget(\cs, \thin b_2, \thin \pick_{5(b_0, b_1)}, \thin \arr[4])}$
$\line{24}{}{\arr'[5] = \pick_{5(b_0, b_1, b_3)}}$

$\line{25}{}{\pick_{\auxb, 6(b_0)}: \CircuitVal \deq \pickgadget(\cs, \thin b_0, \thin \arr[6], \thin \value)}$
$\line{26}{}{\pick_{\auxb, 6(b_0, b_1)}: \CircuitVal \deq \pickgadget(\cs, \thin b_1, \thin \pick_{6(b_0)}, \thin \arr[5])}$
$\line{27}{}{\pick_{\auxb, 6(b_0, b_1, b_2)}: \CircuitVal \deq \pickgadget(\cs, \thin b_2, \thin \pick_{6(b_0, b_1)}, \thin \arr[5])}$
$\line{28}{}{\arr'[6] = \pick_{6(b_0, b_1, b_3)}}$

$\line{27}{}{\pick_{\auxb, 7(b_0, b_1)}: \CircuitVal \deq \pickgadget(\cs, \thin \and_{(b_0, b_1)}, \thin \value, \thin \arr[6])}$
$\line{28}{}{\pick_{\auxb, 7(b_0, b_1, b_2)}: \CircuitVal \deq \pickgadget(\cs, \thin b_2, \thin \pick_{7(b_0, b_1)}, \thin \arr[6])}$
$\line{29}{}{\arr'[7] = \pick_{7(b_0, b_1, b_3)}}$

$\line{30}{}{\return \arr'}$

### AND Gadget

The function $\andgadget$ returns an allocated bit $1$ if both allocated bit arguments $\x$ and $\y$ are $1$ and returns the allocated bit $0$ otherwise. 

**Implementation:** [`bellman::gadgets::boolean::AllocatedBit::and()`](https://github.com/filecoin-project/bellman/blob/e0ac6b879eac87832ba6ef2e37320e877d1d96b6/src/gadgets/boolean.rs#L155)

The $\RCS$ quadratic constraint that is added by the $\andgadget$, when applied to two boolean constrained values $\x, \y \in \Bit$ and outputting a third boolean constrained value $\and \in \Bit$, is:

$\bi (\x) * (\y) = (\and)$

This table shows the satisfiablilty of the constraint for all values of $\x, \y \in \Bit$ and corresponding outputted values of $\and \in \Bit$.

| $\x$   | $\y$    | $\and$ | $(\x) * (\y) = (\and)$ |
|--------|---------|--------|-----------------------|
| $0$    | $0$     | $0$      | $(0) * (0) = (0)$ |
| $1$    | $0$     | $0$      | $(1) * (0) = (0)$ |
| $0$    | $1$     | $0$      | $(0) * (1) = (0)$ |
| $1$    | $1$     | $1$      | $(1) * (1) = (1)$ |

<br />

$\overline{\Function \andgadget(\qquad}$
$\quad \cs: \RCS,$
$\quad \x_\aap: \CircuitBit,$
$\quad \y_\aap: \CircuitBit,$
$\underline{) \rightarrow \CircuitBit \qquad\qquad\bi}$
$\line{1}{\bi}{\and: \Bit = \x_\aap\dot\int \thin\AND\thin \y_\aap\dot\int}$
$\line{2}{\bi}{\and_\auxb: \CircuitBit \deq \cs\dot\privateinput(\and)}$
$\line{3}{\bi}{\cs\dot\assert(\x_\aap * \y_\aap = \and_\auxb)}$
$\line{4}{\bi}{\return \and_\auxb}$

### NOR Gadget

The function $\norgadget$ returns an allocated bit $1$ if both allocated bit arguments $\x$ and $\y$ are $0$ and returns the allocated bit $0$ otherwise. 

**Implementation:** [`bellman::gadgets::boolean::AllocatedBit::nor()`](https://github.com/filecoin-project/bellman/blob/e0ac6b879eac87832ba6ef2e37320e877d1d96b6/src/gadgets/boolean.rs#L231)

The $\RCS$ quadratic constraint that is added by $\norgadget$, when applied to two boolean constrained values $\x, \y \in \Bit$ and outputting a third boolean constrained value $\nor \in \Bit$, is:

$\bi (1 - \x) * (1 - \y) = (\nor)$

The following table shows the satisfiablilty of the constraint for all values of $\x, \y \in \Bit$ and corresponding outputted values for $\nor \in \Bit$.

| $\x$   | $\y$    | $\nor$ | $(1 - \x) * (1 - \y) = (\nor)$ |
|--------|---------|--------|-----------------------|
| $0$    | $0$     | $1$      | $(1) * (1) = (1)$ |
| $1$    | $0$     | $0$      | $(0) * (1) = (0)$ |
| $0$    | $1$     | $0$      | $(1) * (0) = (0)$ |
| $1$    | $1$     | $0$      | $(0) * (0) = (0)$ |

<br />

$\overline{\Function \norgadget(\qquad}$
$\quad \cs: \RCS,$
$\quad \x_\aap: \CircuitBit,$
$\quad \y_\aap: \CircuitBit,$
$\underline{) \rightarrow \CircuitBit \qquad\qquad\bi}$
$\line{1}{\bi}{\nor: \Bit = \neg (\x_\aap\dot\int \OR \y_\aap\dot\int)}$
$\line{2}{\bi}{\nor_\auxb: \CircuitBit \deq \cs\dot\privateinput(\nor)}$
$\line{3}{\bi}{\lc_A: \LinearCombination \equiv 1 - \x_\aap}$
$\line{4}{\bi}{\lc_B: \LinearCombination \equiv 1 - \y_\aap}$
$\line{5}{\bi}{\lc_C: \LinearCombination \equiv \nor_\auxb}$
$\line{6}{\bi}{\cs\dot\assert(\lc_A * \lc_B = \lc_C)}$
$\line{7}{\bi}{\return \nor_\auxb}$

