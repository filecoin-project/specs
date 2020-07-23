---
title: "Notation, Constants, and Types"
weight: 2
math-mode: true
description: "Notation, Constants, and Types for Stacked DRG PoRep"
---

{{< plain hidden >}}
$$
\gdef\createporepbatch{\textsf{create\_porep\_batch}}
\gdef\GrothProof{\textsf{Groth16Proof}}
\gdef\Groth{\textsf{Groth16}}
\gdef\GrothEvaluationKey{\textsf{Groth16EvaluationKey}}
\gdef\GrothVerificationKey{\textsf{Groth16VerificationKey}}
\gdef\creategrothproof{\textsf{create\_groth16\_proof}}
\gdef\ParentLabels{\textsf{ParentLabels}}
\gdef\or#1#2{\langle #1 | #2 \rangle}
\gdef\porepreplicas{\textsf{porep\_replicas}}
\gdef\postreplicas{\textsf{post\_replicas}}
\gdef\winningpartitions{\textsf{winning\_partitions}}
\gdef\windowpartitions{\textsf{window\_partitions}}
\gdef\lebinrep#1{{\llcorner #1 \lrcorner_{2, \textsf{le}}}}
\gdef\bebinrep#1{{\llcorner #1 \lrcorner_{2, \textsf{be}}}}
\gdef\lebytesbinrep#1{{\llcorner #1 \lrcorner_{2, \textsf{le-bytes}}}}
\gdef\fesitelrounds{\textsf{fesitel\_rounds}}
\gdef\int{\textsf{int}}
\gdef\lebytes{\textsf{le-bytes}}
\gdef\lebytestolebits{\textsf{le\_bytes\_to\_le\_bits}}
\gdef\lebitstolebytes{\textsf{le\_bits\_to\_le\_bytes}}
\gdef\letooctet{\textsf{le\_to\_octet}}
\gdef\byte{\textsf{byte}}
\gdef\postpartitions{\textsf{post\_partitions}}
\gdef\PostReplica{\textsf{PostReplica}}
\gdef\PostReplicas{\textsf{PostReplicas}}
\gdef\PostPartitionProof{\textsf{PostPartitionProof}}
\gdef\PostReplicaProof{\textsf{PostReplicaProof}}
\gdef\TreeRProofs{\textsf{TreeRProofs}}
\gdef\pad{\textsf{pad}}
\gdef\octettole{\textsf{octet\_to\_le}}
\gdef\packed{\textsf{packed}}
\gdef\val{\textsf{val}}
\gdef\bits{\textsf{bits}}
\gdef\partitions{\textsf{partitions}}
\gdef\Batch{\textsf{Batch}}
\gdef\batch{\textsf{batch}}
\gdef\postbatch{\textsf{post\_batch}}
\gdef\postchallenges{\textsf{post\_challenges}}
\gdef\Nonce{\textsf{Nonce}}
\gdef\createvanillaporepproof{\textsf{create\_vanilla\_porep\_proof}}
\gdef\PorepVersion{\textsf{PorepVersion}}
\gdef\bedecode{\textsf{be\_decode}}
\gdef\OR{\mathbin{|}}
\gdef\indexbits{\textsf{index\_bits}}
\gdef\nor{\textsf{nor}}
\gdef\and{\textsf{and}}
\gdef\norgadget{\textsf{nor\_gadget}}
\gdef\andgadget{\textsf{and\_gadget}}
\gdef\el{\textsf{el}}
\gdef\arr{\textsf{arr}}
\gdef\pickgadget{\textsf{pick\_gadget}}
\gdef\pick{\textsf{pick}}
\gdef\int{\textsf{int}}
\gdef\x{\textsf{x}}
\gdef\y{\textsf{y}}
\gdef\aap{{\langle \auxb | \pubb \rangle}}
\gdef\aapc{{\langle \auxb | \pubb | \constb \rangle}}
\gdef\TreeRProofs{\textsf{TreeRProofs}}
\gdef\parentlabelsbits{\textsf{parent\_labels\_bits}}
\gdef\label{\textsf{label}}
\gdef\layerbits{\textsf{layer\_bits}}
\gdef\labelbits{\textsf{label\_bits}}
\gdef\digestbits{\textsf{digest\_bits}}
\gdef\node{\textsf{node}}
\gdef\layerindex{\textsf{layer\_index}}
\gdef\be{\textsf{be}}
\gdef\octet{\textsf{octet}}
\gdef\reverse{\textsf{reverse}}
\gdef\LSBit{\textsf{LSBit}}
\gdef\MSBit{\textsf{MSBit}}
\gdef\LSByte{\textsf{LSByte}}
\gdef\MSByte{\textsf{MSByte}}
\gdef\PorepPartitionProof{\textsf{PorepPartitionProof}}
\gdef\PostPartitionProof{\textsf{PostPartitionProof}}
\gdef\octetbinrep#1{{\llcorner #1 \lrcorner_{\lower{2pt}{2, \textsf{octet}}}}}
\gdef\fieldelement{\textsf{field\_element}}
\gdef\Fqsafe{{\mathbb{F}_{q, \safe}}}
\gdef\elem{\textsf{elem}}
\gdef\challenge{\textsf{challenge}}
\gdef\challengeindex{\textsf{challenge\_index}}
\gdef\uniquechallengeindex{\textsf{unique\_challenge\_index}}
\gdef\replicaindex{\textsf{replica\_index}}
\gdef\uniquereplicaindex{\textsf{unique\_replica\_index}}
\gdef\nreplicas{\textsf{n\_replicas}}
\gdef\unique{\textsf{unique}}
\gdef\R{\mathcal{R}}
\gdef\getpostchallenge{\textsf{get\_post\_challenge}}
\gdef\verifyvanillapostproof{\textsf{verify\_vanilla\_post\_proof}}
\gdef\BinPathElement{\textsf{BinPathElement}}
\gdef\BinTreeDepth{\textsf{BinTreeDepth}}
\gdef\BinTree{\textsf{BinTree}}
\gdef\BinTreeProof{\textsf{BinTreeProof}}
\gdef\bintreeproofisvalid{\textsf{bintree\_proof\_is\_valid}}
\gdef\Bit{{\{0, 1\}}}
\gdef\Byte{\mathbb{B}}
\gdef\calculatebintreechallenge{\textsf{calculate\_bintree\_challenge}}
\gdef\calculateocttreechallenge{\textsf{calculate\_octtree\_challenge}}
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
\gdef\octtreeproofisvalid{\textsf{octtree\_proof\_is\_valid}}
\gdef\path{\textsf{path}}
\gdef\pathelem{\textsf{path\_elem}}
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
\gdef\nodeindex{\textsf{node\_index}}
\gdef\childindex{\textsf{child\_index}}
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
\gdef\bintreerootgadget{\textsf{bintree\_root\_gadget}}
\gdef\octtreerootgadget{\textsf{octtree\_root\_gadget}}
\gdef\cs{\textsf{cs}}
\gdef\RCS{\textsf{R1CS}}
\gdef\pathbits{\textsf{path\_bits}}
\gdef\missingbit{\textsf{missing\_bit}}
\gdef\missingbits{\textsf{missing\_bits}}
\gdef\pubb{\textbf{pub}}
\gdef\privb{\textbf{priv}}
\gdef\auxb{\textbf{aux}}
\gdef\constb{\textbf{const}}
\gdef\CircuitVal{\textsf{CircuitVal}}
\gdef\CircuitBit{{\textsf{CircuitVal}_\Bit}}
\gdef\Le{\textsf{le}}
\gdef\privateinput{\textsf{private\_input}}
\gdef\publicinput{\textsf{public\_input}}
\gdef\deq{\mathbin{\overset{\diamond}{=}}}
\gdef\alloc{\textsf{alloc}}
\gdef\insertgadget#1{\textsf{insert\_#1\_gadget}}
\gdef\block{\textsf{block}}
\gdef\shagadget#1#2{\textsf{sha#1\_#2\_gadget}}
\gdef\poseidongadget#1{\textsf{poseidon\_#1\_gadget}}
\gdef\refeq{\mathbin{\overset{{\small \&}}=}}
\gdef\ptreq{\mathbin{\overset{{\small \&}}=}}
\gdef\bit{\textsf{bit}}
\gdef\extend{\textsf{extend}}
\gdef\auxle{{[\textbf{aux}, \textsf{le}]}}
\gdef\SpecificNotation{{\underline{\text{Specific Notation}}}}
\gdef\repeat{\textsf{repeat}}
\gdef\preimage{\textsf{preimage}}
\gdef\digest{\textsf{digest}}
\gdef\digestbytes{\textsf{digest\_bytes}}
\gdef\digestint{\textsf{digest\_int}}
\gdef\leencode{\textsf{le\_encode}}
\gdef\ledecode{\textsf{le\_decode}}
\gdef\ReplicaID{\textsf{ReplicaID}}
\gdef\replicaid{\textsf{replica\_id}}
\gdef\replicaidbits{\textsf{replica\_id\_bits}}
\gdef\replicaidblock{\textsf{replica\_id\_block}}
\gdef\cc{\textsf{::}}
\gdef\new{\textsf{new}}
\gdef\lebitsgadget{\textsf{le\_bits\_gadget}}
\gdef\CircuitBitOrConst{{\textsf{CircuitValOrConst}_\Bit}}
\gdef\createporepcircuit{\textsf{create\_porep\_circuit}}
\gdef\CommD{\textsf{CommD}}
\gdef\CommC{\textsf{CommC}}
\gdef\CommR{\textsf{CommR}}
\gdef\CommCR{\textsf{CommCR}}
\gdef\commd{\textsf{comm\_d}}
\gdef\commc{\textsf{comm\_c}}
\gdef\commr{\textsf{comm\_r}}
\gdef\commcr{\textsf{comm\_cr}}
\gdef\assert{\textsf{assert}}
\gdef\asserteq{\textsf{assert\_eq}}
\gdef\TreeDProof{\textsf{TreeDProof}}
\gdef\TreeRProof{\textsf{TreeRProof}}
\gdef\TreeR{\textsf{TreeR}}
\gdef\ParentColumnProofs{\textsf{ParentColumnProofs}}
\gdef\challengebits{\textsf{challenge\_bits}}
\gdef\packedchallenge{\textsf{packed\_challenge}}
\gdef\PartitionProof{\textsf{PartitionProof}}
\gdef\u#1{\textsf{u#1}}
\gdef\packbitsasinputgadget{\textsf{pack\_bits\_as\_input\_gadget}}
\gdef\treedleaf{\textsf{tree\_d\_leaf}}
\gdef\treerleaf{\textsf{tree\_r\_leaf}}
\gdef\calculatedtreedroot{\textsf{calculated\_tree\_d\_root}}
\gdef\calculatedtreerleaf{\textsf{calculated\_tree\_r\_leaf}}
\gdef\calculatedcommd{\textsf{calculated\_comm\_d}}
\gdef\calculatedcommc{\textsf{calculated\_comm\_c}}
\gdef\calculatedcommr{\textsf{calculated\_comm\_r}}
\gdef\calculatedcommcr{\textsf{calculated\_comm\_cr}}
\gdef\layers{\textsf{layers}}
\gdef\total{\textsf{total}}
\gdef\column{\textsf{column}}
\gdef\parentcolumns{\textsf{parent\_columns}}
\gdef\columns{\textsf{columns}}
\gdef\parentlabel{\textsf{parent\_label}}
\gdef\label{\textsf{label}}
\gdef\calculatedtreecleaf{\textsf{calculated\_tree\_c\_leaf}}
\gdef\calculatedcolumn{\textsf{calculated\_column}}
\gdef\parentlabels{\textsf{parent\_labels}}
\gdef\drg{\textsf{drg}}
\gdef\exp{\textsf{exp}}
\gdef\parentlabelbits{\textsf{parent\_label\_bits}}
\gdef\parentlabelblock{\textsf{parent\_label\_block}}
\gdef\Bits{\textsf{ Bits}}
\gdef\safe{\textsf{safe}}
\gdef\calculatedlabel{\textsf{calculated\_label}}
\gdef\createlabelgadget{\textsf{create\_label\_gadget}}
\gdef\encodingkey{\textsf{encoding\_key}}
\gdef\encodegadget{\textsf{encode\_gadget}}
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
\gdef\merkleproofs{\textsf{merkle\_proofs}}
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
\gdef\V{\mathcal{V}}
\gdef\ww{{\textsf{winning}|\textsf{window}}}
\gdef\replicasperk{{\textsf{replicas}/k}}
\gdef\replicas{\textsf{replicas}}
\gdef\Replica{\textsf{Replica}}
\gdef\createvanillapostproof{\textsf{create\_vanilla\_post\_proof}}
\gdef\createpostcircuit{\textsf{create\_post\_circuit}}
\gdef\ReplicaProof{\textsf{ReplicaProof}}
\gdef\aww{{\langle \ww \rangle}}
\gdef\partitionproof{\textsf{partition\_proof}}
\gdef\replicas{\textsf{replicas}}
\gdef\getdrgparents{\textsf{get\_drg\_parents}}
\gdef\getexpparents{\textsf{get\_exp\_parents}}
\gdef\DrgSeed{\textsf{DrgSeed}}
\gdef\DrgSeedPrefix{\textsf{DrgSeedPrefix}}
\gdef\FeistelKeysBytes{\textsf{FeistelKeysBytes}}
\gdef\porep{\textsf{porep}}
\gdef\rng{\textsf{rng}}
\gdef\ChaCha#1{\textsf{ChaCha#1}}
\gdef\cc{\textsf{::}}
\gdef\fromseed{\textsf{from\_seed}}
\gdef\buckets{\textsf{buckets}}
\gdef\meta{\textsf{meta}}
\gdef\dist{\textsf{dist}}
\gdef\each{\textsf{each}}
\gdef\PorepID{\textsf{PorepID}}
\gdef\porepgraphseed{\textsf{porep\_graph\_seed}}
\gdef\utf{\textsf{utf8}}
\gdef\DrgStringID{\textsf{DrgStringID}}
\gdef\FeistelStringID{\textsf{FeistelStringID}}
\gdef\graphid{\textsf{graph\_id}}
\gdef\createfeistelkeys{\textsf{create\_feistel\_keys}}
\gdef\FeistelKeys{\textsf{FeistelKeys}}
\gdef\feistelrounds{\textsf{fesitel\_rounds}}
\gdef\feistel{\textsf{feistel}}
\gdef\ExpEdgeIndex{\textsf{ExpEdgeIndex}}
\gdef\loop{\textsf{loop}}
\gdef\right{\textsf{right}}
\gdef\left{\textsf{left}}
\gdef\mask{\textsf{mask}}
\gdef\RightMask{\textsf{RightMask}}
\gdef\LeftMask{\textsf{LeftMask}}
\gdef\roundkey{\textsf{round\_key}}
\gdef\beencode{\textsf{be\_encode}}
\gdef\Blake{\textsf{Blake2b}}
\gdef\input{\textsf{input}}
\gdef\output{\textsf{output}}
\gdef\while{\textsf{while }}
\gdef\digestright{\textsf{digest\_right}}
\gdef\xor{\mathbin{\oplus_\text{xor}}}
\gdef\Edges{\textsf{ Edges}}
\gdef\edge{\textsf{edge}}
\gdef\expedge{\textsf{exp\_edge}}
\gdef\expedges{\textsf{exp\_edges}}
\gdef\createlabel{\textsf{create\_label}}
\gdef\Label{\textsf{Label}}
\gdef\Column{\textsf{Column}}
\gdef\Columns{\textsf{Columns}}
\gdef\ParentColumns{\textsf{ParentColumns}}
% `\tern` should be written as
% \gdef\tern#1?#2:#3{#1\ \text{?}\ #2 \ \text{:}\ #3}
% but that's not possible due to https://github.com/KaTeX/KaTeX/issues/2288
\gdef\tern#1#2#3{#1\ \text{?}\ #2 \ \text{:}\ #3}
\gdef\repeattolength{\textsf{repeat\_to\_length}}
\gdef\verifyvanillaporepproof{\textsf{verify\_vanilla\_porep\_proof}}
\gdef\poreppartitions{\textsf{porep\_partitions}}
\gdef\challengeindex{\textsf{challenge\_index}}
\gdef\porepbatch{\textsf{porep\_batch}}
\gdef\winningchallenges{\textsf{winning\_challenges}}
\gdef\windowchallenges{\textsf{window\_challenges}}
\gdef\PorepPartitionProof{\textsf{PorepPartitionProof}}
\gdef\TreeD{\textsf{TreeD}}
\gdef\TreeCProof{\textsf{TreeCProof}}
\gdef\Labels{\textsf{Labels}}
\gdef\porepchallenges{\textsf{porep\_challenges}}
\gdef\postchallenges{\textsf{post\_challenges}}
\gdef\PorepChallengeSeed{\textsf{PorepChallengeSeed}}
\gdef\getporepchallenges{\textsf{get\_porep\_challenges}}
\gdef\getallparents{\textsf{get\_all\_parents}}
\gdef\PorepChallengeProof{\textsf{PorepChallengeProof}}
\gdef\challengeproof{\textsf{challenge\_proof}}
\gdef\PorepChallenges{\textsf{PorepChallenges}}
\gdef\replicate{\textsf{replicate}}
\gdef\createreplicaid{\textsf{create\_replica\_id}}
\gdef\ProverID{\textsf{ProverID}}
\gdef\replicaid{\textsf{replica\_id}}
\gdef\generatelabels{\textsf{generate\_labels}}
\gdef\labelwithdrgparents{\textsf{label\_with\_drg\_parents}}
\gdef\labelwithallparents{\textsf{label\_with\_all\_parents}}
\gdef\createtreecfromlabels{\textsf{create\_tree\_c\_from\_labels}}
\gdef\ColumnDigest{\textsf{ColumnDigest}}
\gdef\encode{\textsf{encode}}
\gdef\sector{\textsf{sector}}
$$
{{< /plain >}}

# SDR Notation, Constants, and Types
---

## General Notation

`$\mathbb{T}^{[n]}$`\
An array of `$n$` elements of type `$\mathbb{T}$`.

`$V[i]$`\
The element of array `$V$` at index `$i$`. **All indexes in this document (unless otherwise noted) start at `$0$`.**

`$[n] \equiv [0, n - 1] = 0, 1, \ldots, n - 1$`\
The range of integers from `$0$` up to, but not including, `$n$`. Despite the use of square brackets `$n$` is not included in this range.

`$\int: [n]$`\
The type `$[n]$` denotes `$\int$` as being an integer in the range `$[n] \equiv 0, \ldots, n - 1 \thin$`.

`$\Byte$`\
The byte type, an integer in `$[256]$`.

`$\Bit^{[n]}$`\
A bit string of length `$n$`.

`$0^{[n]}$`\
`$1^{[n]}$`\
Creates a bit-string where each of the `$n$` bits is set to `$0$` or `$1$` respectively.

`$\u{32} = \Bit^{[32]}$`\
`$\u{64} = \Bit^{[64]}$`\
32 and 64-bit unsigned integers equipped with standard arithmetic operations.

`$\u{64}_{(17)}$`\
A 64-bit unsigned integer where only the least significant 17 bits are utilized (can be `$0$` or `$1$`, the remaining unutilized 47 bits are unused and set to `$0$`).

`$\Fq \equiv [q] \equiv 0, \ldots, q - 1$`\
A field element. An element of curve BLS12-381's scalar field, an unsigned integer in `$[q]$` equipped with the field's binary operations. The scalar field has prime order `$q$`.

`$[a, b] = a, a + 1, \ldots, b$`\
The range of integers from `$a$` up to and including `$b$` (both endpoints are inclusive). There is an **ambiguity in notation** between this and the construction of a two element array. The notation `$[a, b]$` **always refers to a range except in the cases where we pass a two element array as an argument into the hash functions `$\Sha{254}_2$` and `$\Poseidon_2$`**.

`$V[a..b] = V_a, \ldots, V_{b - 1}$`\
The slice of elements of array `$V$` from index `$a$` (**inclusive**) up to index `$b$` (**noninclusive**).

`$V[..n] \equiv V[0 .. n] = V_0, \ldots, V_{n - 1}$`\
A slice of `$V$` starting at the first element (**inclusive**) and ending at the element at index `$n$` (**noninclusive**).

`$V[n..] \equiv V[n..\len(V)] = V_n, \dots, V_{\len(V) - 1}$`\
A slice of `$V$` containing the elements starting at index `$n$` (inclusive) up to and including the last element.

`$M[r][:]$`\
The `$r^{th}$` row of matrix `$M$` as an array.

`$M[:][c] = [\thin M[r][c]\ |\ \forall r \in [n_\text{rows}] \thin]$`\
The `$c^{th}$` column of matrix `$M$` (having `$n_\text{rows}$` rows) as a flattened array.

`$\mathbb{T}^{[n]} \concat \mathbb{T}^{[m]} \rightarrow \mathbb{T}^{[n + m]}$`\
Concatenates two arrays (whose elements are the same type) producing an array of type `$\mathbb{T}^{[n + m]} \thin$`.

`$a \MOD b$`\
Integer `$a$` modulo integer `$b$`. Maps `$a$` into the range `$0, \ldots, b - 1 \thin$`.

`$x \leftarrow\ S$`\
Samples `$x$` uniformly from `$S$`.

`$a \ll n$`\
Bitwise left shift of `$a$` by `$n$` bits. **The leftmost bit of a unsigned integer is the most significant**, for example `$\int: \u{8} = 5_{10} = 00000101_2 \thin$`.

`$a \gg n$`\
Bitwise right shift of `$a$` by `$n$` bits. **The rightmost bit of a unsigned integer is the least significant**, for example `$(\int \gg 1) \AND 1$` returns the 2nd least significant bit for an integer `$\int: \u{32}, \u{64} \thin$`.

`$a \AND b:$`\
Bitwise AND of `$a$` and `$b$`, where `$a$` and `$b$` are integers or bit-strings. **The rightmost bit of a unsigned integer is the least significant**, for example `$\int \AND 1$` returns the least significant bit for an integer `$\int: \u{32}, \u{64} \thin$`.

`$a \OR b$`\
Bitwise OR. Used in conjunction with bitwise shift to concatenate bit strings, e.g. `$(101_2 \ll 3) \OR 10_2 = 10110_2 \thin$`.

`$a \xor b$`\
Bitwise XOR.

`$a \oplus b$`\
`$a \ominus b$`\
Field addition and subtraction respectively. `$a$` and `$b$` are field elements `$\Fq$`.

`$a, b \Leftarrow X$`\
Destructures and instance of `$\struct X\ \{a: \mathbb{A},\ b: \mathbb{B}\}$` into two variables `$a$` and `$b$`, assigning to each variable the value (and type) of `$X.a$` and `$X.b$` respectively.

`$\for \mathbb{T}\ \{\ x ,\thin y \ \} \in \textsf{Iterator}:$`\
Iterates over `$\textsf{Iterator}$` and destructures each element into variables `$x$` and `$y$` local within the for-loop iteration's scope. Each element of `$\textsf{Iterator}$` is an instance of a structure `$\mathbb{T}$`.

`$\sum_{i \in [n]}{\textsf{expr}_i}$`\
For each `$i \in 0, \ldots, n - 1 \thin$`, summates the value output by the `$i^{th}$` expression `$\textsf{expr}_i\thin$`.

`$\big\|_{x \hspace{1pt} \in \hspace{1pt} S} \thin x$`\
The concatenation of all elements `$x$` of an iterator `$S$`.

`$V = [x\ \textsf{as}\ \mathbb{T}\ |\ \forall x \in S]$`\
Array builder notation. The above creates an array `$V$` of type `$\mathbb{T}^{[\textsf{len}(S)]}$` where each `$V[i] = S[i] \as \mathbb{T}$` .

`$\llcorner \int \lrcorner_{2, \Le}: \Bit^{[\lceil \log_2(\int) \rceil]}$`\
`$\llcorner \int \lrcorner_{8, \Le}: [8]^{[\lceil \log_8(\int) \rceil]}$`\
The little-endian binary and octal representations of an unsigned integer `$\int$`.

`$\for \each \in \textsf{Iterator}:$`\
Iterates over each element of `$\textsf{Iterator}$` and ignores each element's value.

`$\big[ \tern{x = 0}{a}{b \big}]$`\
A ternary expression. Returns `$a$` if `$x = 0$`, otherwise returns `$b$`.

`$\mathbb{T}^{[m]}\dot\extend(\mathbb{T}^{[n]})  \rightarrow \mathbb{T}^{[m + n]}$`\
Appends each value of `$\mathbb{T}^{[n]}$` onto the end of the array `$\mathbb{T}^{[m]}$`. Equivalent to writing:

```text
$\bi \line{1}{\bi}{x: \mathbb{T}^{[m]} = [x_1, \thin\ldots\thin, x_m]}$
$\bi \line{2}{\bi}{y: \mathbb{T}^{[n]} = [y_1, \thin\ldots\thin, y_n]}$
$\bi \line{3}{\bi}{x: \mathbb{T}^{[m + n]} = x \concat y}$
```

`$\mathbb{T}^{[m]}\dot\repeattolength(n: \mathbb{N} \geq m) \rightarrow \mathbb{T}^{[n]}$`\
Repeats the values of array `$\mathbb{T}^{[m]}$` until its length is `$n$`. Equivalent to writing:

```text
$\line{1}{\bi}{x: \mathbb{T}^{[m]} = [x_1, \thin\ldots\thin, x_m]}$
$\line{2}{\bi}{\for i \in [n - m]:}$
$\line{3}{\bi}{\quad x\dot\push(x[i \MOD m])}$
```


## Protocol Constants

**Implementation:** [`filecoin_proofs::constants`](https://github.com/filecoin-project/rust-fil-proofs/blob/d088e7b997c13a59e66c062a9ceb110f991c9849/filecoin-proofs/src/constants.rs)

`$\ell_\sector^\byte = 32\ \textsf{GiB} = 32 * 1024 * 1024 * 1024\ \textsf{Bytes}$`\
The byte length of a sector `$D$`.

`$\ell_\node^\byte = \ell_\Fq^\byte = 32\ \textsf{Bytes}$`\
The byte length of a node.

`$N_\nodes = \ell_\sector^\byte / \ell_\node^\byte = 2^{30}\ \textsf{Nodes}$`\
The number of nodes in each DRG. The protocol guarantees that the sector byte length is evenly divisible by the node byte length.

`$\ell_\Fq^\bit = \lceil \log_2(q) \rceil = \lceil 254.85\ldots \rceil = 255 \Bits$`\
The number of bits required to represent any field element.

`$\ell_\Fqsafe^\bit = \ell_\Fq^\bit - 1 = 254 \Bits$`\
The maximum integer number of bits that can be safely casted (casted without the possibility of failure) into a field element.

`$\ell_\block^\bit = 256 \Bits$`\
The bit length of a `$\Sha{256}$` block.

`$d_\exp = 6$`\
The degree of each DRG. The number of DRG parents generated for each node. The total number of parents generated for nodes in the first Stacked DRG layer.

`$d_\meta = d_\drg - 1 = 5$`\
The degree of the DRG metagraph. The number of DRG parents generated using a metagraph.

`$d_\exp = 8$`\
The degree of each expander. The number of expander parents generated for each node not in the first Stacked-DRG layer.

`$d_\total = d_\drg + d_\exp = 14$`\
The total number of parents generated for all nodes not in the first Stacked-DRG layer.

`$N_\expedges = d_\exp * N_\nodes  = 2^{33} \Edges$`\
The number of edges per expander graph.

`$\ell_\expedge^\bit = \log_2(N_\expedges) = 33 \Bits$`\
The number of bits required to represent the index of an expander edge.

`$\ell_\mask^\bit = \lceil \ell_\expedge^\bit / 2 \rceil = 17 \Bits$`\
The number of bits that are masked by a Feistel network bitmask.

`$\RightMask: \u{64}_{(17)} = 2^{\ell_\mask^\bit} - 1 \qquad\quad\quad\bi\ = 0000000000000000011111111111111111_2$`\
`$\LeftMask: \u{64}_{(17 \ldotdot 34)} = \RightMask \ll \ell_\mask^\bit = 1111111111111111100000000000000000_2$`\
The Fesitel network's right-half and left-half bitmasks. Each bitmask contains `$\ell_\mask^\bit$` bits set to `$1$`. Both bitmasks are represented in binary as `$34 = 2 * \ell_\mask^\bit$` digits. Note that `$\RightMask$`'s lowest 17 bits are utilized and `$\LeftMask$`'s lowest 17th-34th bits are utilized.

`$N_\feistelrounds = 4$`\
The number of rounds per Feistel network.

`$N_\layers = 11$`\
The number of DRG layers in the Stacked DRG.

`$N_\buckets = \lceil \log_2(N_\nodes) \rceil = 30$`\
The number of buckets used in the Bucket Sample algorithm to generate DRG parents.

`$N_\parentlabels = 37$`\
The number of parent labels factored into each node label.

`$\BinTreeDepth = \log_2(N_\nodes) = 30$`\
`$\OctTreeDepth = \log_8(N_\nodes) = 10$`\
The depth of a `$\BinTree$` and `$\OctTree$` respectively. The number of tree layers is the tree's depth `$+ 1 \thin$`. The Merkle hash arity of trees are 2 and 8 respectively.

`$q = \text{73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001}_{16}$`\
`$q = 52435875175126190479447740508185965837690552500527637822603658699938581184513_{10}$`\
The (prime) order of the `$\textsf{BLS12-381}$` scalar field (given in hex and decimal).

`$N_{\poreppartitions / \batch} = 10$`\
`$N_{\postpartitions / \batch, \P, \thin \winning} \leq \len(\PostReplicas_{\P, \batch})$`\
`$N_{\postpartitions / \batch, \P, \thin \window} \leq \len(\PostReplicas_{\P, \batch})$`\
The number of partition proofs per PoRep, Winning PoSt, and Window PoSt proof batch. The number of PoSt partition proofs in a batch is specific to the size of the PoSt prover `$\P$`'s replica set `$\PostReplicas_{\P, \batch}$` at the time of batch proof generation.

`$N_{\porepreplicas / k} = 1$`\
`$N_{\postreplicas / k, \winning} = 1$`\
`$N_{\postreplicas / k, \window} = 2349$`\
The number of challenged replicas per PoRep, Winning PoSt, and Window PoSt partition proof.

`$N_{\postreplicas / k \thin \aww} \in \{ N_{\postreplicas / k \thin, \winning}, N_{\postreplicas / k, \window} \}$`\
Notational shortand meaning "either Winning of Window PoSt, determined by context".

`$N_{\porepchallenges / k} \equiv N_{\porepchallenges / R} = 176$`\
The number of Merkle challenges per PoRep partition proof. PoRep partition proofs are generated using a single replica `$R$`.

`$N_{\postchallenges / R, \thin \winning} = 66$`\
`$N_{\postchallenges / R, \thin \window} = 10$`\
The number of Merkle challenges per challenged replica `$R$` in a PoSt partition proof.

`$\GrothEvaluationKey_{\langle \textsf{circ} \rangle} = \text{<set during}\ \textsf{circ} \text{'s trusted setup>}$`\
`$\GrothVerificationKey_{\langle \textsf{circ} \rangle} = \text{<set during}\ \textsf{circ} \text{'s trusted setup>}$`\
The Groth16 keypair sued to generate and verify SNARKs for a circuit definition `$\textsf{circ}$` (PoRep, Winning PoSt, Window PoSt each for a given protocol sector size).

`$\DrgStringID: {\Byte_\utf}^{[*]} = ``\text{Filecoin\_DRSample}"$`\
`$\FeistelStringID: {\Byte_\utf}^{[*]} = ``\text{Filecoin\_Feistel}"$`\
The ID strings associated with the DRG and Feistel network.

`$\DrgSeedPrefix_\PorepID: \Byte^{[28]} = \Sha{256}(\DrgStringID \concat \PorepID)[\ldotdot 28]$`\
Part of the RNG seed used to sample DRG parents for every PoRep having the version `$\PorepID$`.
[`storage_proofs::core::crypto::derive_porep_domain_seed()`](https://github.com/filecoin-project/rust-fil-proofs/blob/b9126bf56cfd2c73ce4ce1f6cb46fe001450f326/storage-proofs/core/src/crypto/mod.rs#L13)

`$\PorepVersion_\textsf{SDR,32GiB,v1}: \u{64} = 3$`\
The version ID of the PoRep in use. PoRep versions are parameterized by the triple: (PoRep proof system, sector-size, version number).

`$\Nonce_\PorepVersion: \u{64} = 0$`\
Each `$\PorepVersion$`'s has an associated nonce used to generate the PoRep versions `$\PorepID$`. Currently all PoRep vefrsion's have a nonce of `$0$`.

`$\FeistelKeysBytes_\PorepID: \Byte^{[32]} = \Sha{256}(\FeistelStringID \concat \PorepID)$`\
The byte array representing the concatenation of the Feistel network's `$N_\feistelrounds = 4$` 64-bit round keys. All PoRep's corresponding to the version `$\PorepID$` use the same Fesistel keys.\
**Implementation:** [`storage_proofs::core::crypto::derive_porep_domain_seed()`](https://github.com/filecoin-project/rust-fil-proofs/blob/b9126bf56cfd2c73ce4ce1f6cb46fe001450f326/storage-proofs/core/src/crypto/mod.rs#L13)

```text
$\FeistelKeys_\PorepID: \u{64}^{[N_\feistelrounds]}\thin =\thin [$
$\quad \ledecode(\FeistelKeysBytes_\PorepID[\ldotdot 8]) \thin,$
$\quad \ledecode(\FeistelKeysBytes_\PorepID[8 \ldotdot 16]) \thin,$
$\quad \ledecode(\FeistelKeysBytes_\PorepID[16 \ldotdot 24]) \thin,$
$\quad \ledecode(\FeistelKeysBytes_\PorepID[24 \ldotdot]) \thin,$
$]$
```
The Feistel round keys used by every PoRep having version `$\PorepID$`.\
**Implementation:** [`storage_proofs::porep::stacked::vanilla::graph::StackedGraph::new()`](https://github.com/filecoin-project/rust-fil-proofs/blob/b6bf96faff1bf7611982ca1318c20d256fd464d4/storage-proofs/porep/src/stacked/vanilla/graph.rs#L189-#L194)

`$\DrgSeed_{\PorepID, v}: \Byte^{[32]} = \DrgSeedPrefix_\PorepID \concat \leencode(v) \as \Byte^{[4]}$`\
The DRG parent sampling RNG's seed for each node `$v \in [N_\nodes]$` for a porep version `$\PorepID \thin$`.

## Protocol Types and Notation

`$\Safe = \Byte^{[32]}[..31] \concat (\Byte^{[32]}[31] \AND 00111111_2)$`\
An array of 32 bytes that utilizes only the first `$\ell_\Fqsafe^\bit = 254$` bits. `$\Safe$` types can be safely casted into field elements (casting is guaranteed not to fail). `$\Safe$` is least significant byte first.

`$\Fqsafe = \Fq \AND 1^{[254]}$`\
A field element that utilizes only its first `$\ell_\Fqsafe^\bit = 254$` bits. `$\Fqsafe$`'s are created from casting a `$\Safe$` into an `$\Fq$` (the casting produces an `$\Fqsafe \thin$`).


`$\leencode(\Fq) \rightarrow \Byte^{[32]}$`\
`$\leencode(\Fqsafe) \rightarrow \Safe$`\
The produced byte array is least significant byte first.

`$\NodeIndex: \u{32} \in [N_\nodes]$`\
The index of a node in a DRG. Node indexes are 32-bit integers.

`$v: \NodeIndex$`\
`$u: \NodeIndex$`\
Represents the node indexes for a child node `$v$` and a parent `$u$` of `$v$`.

`$\mathbf{u}_\drg: \NodeIndex^{[d_\drg]} = \getdrgparents(v)$`\
An array DRG parents for a child node `$v: \NodeIndex$`. DRG parents are in `$v$`'s Stacked-DRG layer.

`$\mathbf{u}_\exp: \NodeIndex^{[d_\exp]} = \getexpparents(v)$`\
An array of expander parents for a child node `$v: \NodeIndex$`. Expander parents are in the Stacked-DRG layer preceding `$v$`'s.

`$\mathbf{u}_\total: \NodeIndex^{[d_\total]} = \getallparents(v) = \mathbf{u}_\drg \concat \mathbf{u}_\exp$`\
An array containing both the DRG and Expander parents for a child node `$v: \NodeIndex$`. The first `$d_\drg$` elements of `$\mathbf{u}_\total$` are `$v$`'s DRG parents and the last `$d_\exp$` elements are `$v$`'s expander parents.

`$p \in [\len(\mathbf{u}_{\langle \drg | \exp | \total \rangle})]$`\
The index of a parent in a parent node array `$\mathbf{u}_{\langle \drg | \exp | \total \rangle}$`. `$p$` is not the same as a parent node's index `$u: \NodeIndex$`.

`$l \in [N_\layers]$`\
The index of a layer in the Stacked DRG. **This document indexes layers a starting at `$0$` not `$1$`**. Within the context of Merkle trees and proofs, `$l$` may denote a tree layer `$l \in [0, \langle \BinTreeDepth | \OctTreeDepth \rangle] \thin$`.

`$l_v \in [N_\layers]$`\
The Stacked-DRG layer that a node `$v$` resides in.

`$k \in [N_\poreppartitions]$`\
`$k \in [N_{\postpartitions \thin \aww}]$`\
The index of a PoRep or PoSt partition.

`$c: \NodeIndex$`\
A Merkle challenge.

`$\R$`\
A random value.

`$\P, \V$`\
A prover and verifier respectively.

`$X^\dagger$`\
The superscript `$^\dagger$` denotes an unverified value, for example when verifying a proof.

### Merkle Trees

Filecoin utilizes two Merkle tree types, a binary tree type `$\BinTree$` and octal tree type `$\OctTree$`.

`$\BinTree{\bf \sf s}$` use the `$\Sha{254}_2$` hash function (two inputs per hash function call, tree nodes have type `$\Safe$`). `$\OctTree{\bf \sf s}$` use the hash function `$\Poseidon_8$` (eight inputs per hash function call, tree nodes have type `$\Fq$`).

**Implementation:**
* [`merkle_light::merkle::MerkleTree`](https://github.com/filecoin-project/merkle_light/blob/64a468807c594d306d12d943dd90cc5f88d0d6b0/src/merkle.rs#L138)
* [`storage_proofs::core::merkle::MerkleTreeWrapper`](https://github.com/filecoin-project/rust-fil-proofs/blob/f75804c503d9b97a2b02ef3ea2e5d44e8e2c6470/storage-proofs/core/src/merkle/tree.rs#L53)
* [`storage_proofs::core::merkle::MerkleTreeTrait`](https://github.com/filecoin-project/rust-fil-proofs/blob/f75804c503d9b97a2b02ef3ea2e5d44e8e2c6470/storage-proofs/core/src/merkle/tree.rs#L19)

```text
$\struct \BinTree\ \{$
$\quad \leaves \sim \layer_0: \Safe^{[N_\nodes]},$
$\quad \layer_1: \Safe^{[N_\nodes / 2]} \thin,$
$\quad \layer_2: \Safe^{[N_\nodes / 2^2]} \thin,$
$\quad \dots\ ,$
$\quad \layer_{\BinTreeDepth - 1}: \Safe^{[2]},$
$\quad \root \sim \layer_\BinTreeDepth: \Safe,$
$\}$
```
A binary Merkle tree with hash function arity-2 (each Merkle hash operates on two input values). The hash function for `$\BinTree$`'s is `$\Sha{254}_2$` (`$\Sha{256}$` that operates on two 32-byte inputs where the last two bits of the last byte of their `$\Sha{256}$` digest have been zeroed to produce a `$\Safe$` digest). The fields `$\layer_0, \ldots, \layer_\BinTreeDepth$` are arrays containing each tree layer's node labels. The fields `$\leaves$` and `$\root$` are aliases for the fields `$\layer_0$` and `$\layer_\BinTreeDepth$` respectively.

```text
$\struct \OctTree\ \{$
$\quad \leaves \sim \layer_0: \Fq^{[N_\nodes]},$
$\quad \layer_1: \Fq^{[N_\nodes / 8]} \thin,$
$\quad \layer_2: \Fq^{[N_\nodes / 8^2]} \thin,$
$\quad \dots\ ,$
$\quad \layer_{\OctTreeDepth - 1}: \Fq^{[8]},$
$\quad \root \sim \layer_\OctTreeDepth: \Fq,$
$\}$
```
An octal Merkle tree with hash function arity-8 (each Merkle hash operates on 8 input values). The hash function for `$\OctTree$`'s is `$\Poseidon_8$`. The fields `$\layer_0, \ldots, \layer_{10}$` are arrays containing each tree layer's node labels. The fields `$\leaves$` and `$\root$` are aliases for the fields `$\layer_0$` and `$\layer_{10}$` respectively.

### Merkle Proofs

**Implementation:**
* `$\BinTreeProof$`, `$\OctTreeProof$`: [`storage_proofs::merkle::SingleProof`](https://github.com/filecoin-project/rust-fil-proofs/blob/f75804c503d9b97a2b02ef3ea2e5d44e8e2c6470/storage-proofs/core/src/merkle/proof.rs#L344)
* `$\BinPathElement$`, `$\OctPathElement$`: [`storage_proofs::merkle::PathElement`](https://github.com/filecoin-project/rust-fil-proofs/blob/f75804c503d9b97a2b02ef3ea2e5d44e8e2c6470/storage-proofs/core/src/merkle/proof.rs#L223)

```text
$\struct \BinTreeProof\ \{$
$\quad \leaf: \Safe \thin,$
$\quad \root: \Safe \thin,$
$\quad \path: \BinPathElement^{[\BinTreeDepth]} \thin,$
$\}$
```
A `$\BinTree$` Merkle proof generated for a challenge. The notation `$\BinTreeProof_c$` denotes a proof generated for a challenge `$c$`. The field `$\path$` contains one element for tree layers `$0, \ldots, \BinTreeDepth - 1$` (there is no `$\path$` element for the root layer). The path element `$\BinTreeProof\dot\path[l]$` for tree layers `$l \in 0, \ldots, \BinTreeDepth - 1$` contains one node label in layer `$l$` that the Merkle proof verifier will use to calculate the label for a node in layer `$l + 1 \thin$`. Each path element is for a distinct tree layer.

```text
$\struct \BinPathElement\ \{$
$\quad \sibling: \Safe \thin,$
$\quad \missing: \Bit \thin,$
$\}$
```
A single element in `$\BinTreeProof\dot\path$` associated with a single `$\BinTree$` tree layer `$l$`.  Contains the information necessary for a Merkle proof verifier to calculate the label for a node in tree layer `$l + 1$`. The field `$\sibling$` contains the label that the Merkle proof verifier did not calculate for layer `$l$`. The Merkle verifier applies the hash function `$\Sha{254}_2$` to an array of two elements in layer `$l$` to produce the label for a node in layer `$l + 1$`. The order of the elements in each Merkle hash function's 2-element array is given by the field `$\missing$`. If `$\missing = 0$`, then `$\sibling$` is at index `$1$` in the Merkle hash inputs array. Conversely, if `$\missing = 1$` the field `$\sibling$` is at index `$0$` in the 2-element Merkle hash inputs array.

```text
$\struct \OctTreeProof\ \{$
$\quad \leaf: \Fq \thin,$
$\quad \root: \Fq \thin,$
$\quad \path: \OctPathElement^{[\OctTreeDepth]} \thin,$
$\}$
```
A `$\OctTree$` Merkle proof generated for a challenge. The notation `$\OctTreeProof_c$` denotes a proof generated for a challenge `$c$`. The field `$\path$` contains one element for tree layers `$0, \ldots, \OctTreeDepth - 1$` (there is no `$\path$` element for the root layer). The path element `$\OctTreeProof\dot\path[l]$` for tree layers `$l \in 0, \ldots, \OctTreeDepth - 1$` contains one node label in layer `$l$` that the Merkle proof verifier will use to calculate the label for a node in layer `$l + 1 \thin$`. Each path element is for a distinct tree layer. If `$l = 0$` the verifier inserts `$\BinTreeProof\dot\leaf$` into the first path element's Merkle hash inputs array at index `$\BinTreeProof\dot\path[0]\dot\missing \thin$`.

```text
$\struct \OctPathElement\ \{$
$\quad \siblings: \Fq^{[7]} \thin,$
$\quad \missing: [8] \thin,$
$\}$
```
A single element in `$\OctTreeProof\dot\path$` associated with a single `$\OctTree$` tree layer `$l$`.  Contains the information necessary for a Merkle proof verifier to calculate the label for a node in tree layer `$l + 1$`. The field `$\sibling$` contains the label that the Merkle proof verifier did not calculate for layer `$l$`. The Merkle verifier applies the hash function `$\Poseidon_8$` to an array of eight elements in layer `$l$` to produce the label for a node in layer `$l + 1$`. The order of the elements in each Merkle hash function's 8-element array is given by the field `$\missing$`. `$\missing$` is an index in an array of 8 elements telling the verifier at which index in the Merkle hash inputs array that the node label calculated by the verifier for this path element's layer is to be inserted into `$\siblings$`. Given an `$\OctPathElement = \OctTreeProof\dot\path[l]$` in tree layer `$l$`, the node label to be inserted into the hash inputs array at index `$\missing$` was calculated using the path element `$\OctTreeProof\dot\path[l]$` (or `$\OctTreeProof\dot\leaf$` if `$l = 0$`).

### Graphs

`$\ExpEdgeIndex: \u{64}_{(33)} \equiv [N_\expedges]$`\
The index of an edge in an expander graph. Note that `$\ell_\expedge^\bit = 33 \thin$`.

### PoRep

```text
$\PorepID: \Byte^{[32]} =$
$\quad \leencode(\PorepVersion) \as \Byte^{[8]}$
$\quad \|\ \leencode(\Nonce_\PorepVersion) \as \Byte^{[8]}$
$\quad \|\ 0^{[16]}$
```
A unique 32-byte ID assigned to each PoRep version. Each PoRep version is defined by a distinct triple of parameters: (PoRep proof system, sector-size, version number). All PoRep's having the same PoRep version triple have the same `$\PorepID$`. The notation `$0^{[16]}$` denotes 16 zero bytes (not bits).\
**Implementation:** [`filecoin_proofs_api::registry::RegisteredSealProof::porep_id()`](https://github.com/filecoin-project/rust-filecoin-proofs-api/blob/2a01ecc2ff2b9a70fa483d7e36ab79d4f035cf60/src/registry.rs#L110)

`$\SectorID: \u{64}$`\
A unique 64-bit ID assigned to each distinct sector `$D$`.

`$\ReplicaID: \Fqsafe$`\
A unique 32-byte ID assigned to each distinct replica `$R$`.\
**Implementation:** [`storage_proofs::porep::stacked::vanilla::params::generate_replica_id()`](https://github.com/filecoin-project/rust-fil-proofs/blob/be90253f8ff7316d8ef862a1aaed92e76c05ce36/storage-proofs/porep/src/stacked/vanilla/params.rs#L736)

`$\R_\replicaid: \Byte^{[32]}$`\
A random value used to generate a replica `$R$`'s `$\ReplicaID$`.

`$\TreeD: \BinTree$`\
`$\TreeC: \OctTree$`\
`$\TreeR: \OctTree$`\
Merkle trees built over a sector `$D$`, the column digests of a `$\Labels$` matrix, and a replica `$R$` respectively. `$\TreeD$` uses the Merkle hash function `$\Sha{254}_2$` while `$\TreeC$` and `$\TreeR$` use the Merkle hash function `$\Poseidon_8$`. The leaves of `$\TreeD$` are the array of `$D_i: \Safe \in D$` for a sector `$D$`. The leaves of `$\TreeC$` are the array of `$\ColumnDigest$`'s for a replica's labeling `$\Labels$`. The leaves of `$\TreeR$` are the array of `$R_i: \Fq \in R$` for a replica `$R$`.

`$\CommD: \Safe = \TreeD\dot\root$`\
`$\CommC: \Fq = \TreeC\dot\root$`\
`$\CommR: \Fq = \TreeR\dot\root$`\
`$\CommCR: \Fq = \Poseidon_2([\CommC, \CommR])$`\
The Merkle roots of a `$\TreeD$`, `$\TreeC$`, and `$\TreeR$` as well as the commitment `$\CommCR$`, a commitment to both a replicas's Stacked-DRG labeling `$\Labels$` and the replica `$R$`.

`$\Label_{v, l}: \Fqsafe$`\
The label of a node `$v$` in a Stacked-DRG layer `$l \thin$`. The label of a node in a Stacked-DRG layer is specific to a replica `$R$`'s labeling `$\Labels \thin$`. Labels are `$\Fqsafe$` because the labeling function is `$\Sha{254}$` which returns 254-bit/`$\Safe$` digests which are converted into 254-bit/safe field elements.

`$\Labels: {\Label^{[N_\nodes]}}^{[N_\layers]}$`\
An `$N_\layers$` x `$N_\nodes$` matrix containing the label of every node in a replica's Stacked-DRG labeling `$\Labels \thin$`.

`$\Column_v: \Label^{[N_\layers]} = \Labels[:][v]$`\
The label of a node `$v$` in every Stacked-DRG layer (first layer's label is at index `$0$` in the column). A node column is specific to the labeling of a single replica.

`$\ColumnDigest_v: \Fq = \Poseidon_{11}(\Column_v)$`\
The digest of a node `$v$`'s column in a replica's Stacked-DRG. `$\ColumnDigest{\bf \sf s}$` are used as the leaves for `$\TreeC$`. The set of `$\ColumnDigest_v{\bf \sf s}$` for all DRG nodes `$v \in [N_\nodes]$` is specific to a single replica's labeling.

```text
$\struct \ColumnProof\ \{$
$\quad \leaf: \Fq,$
$\quad \root: \Fq,$
$\quad \path: \OctPathElement^{[\OctTreeDepth]},$
$\quad \column: \Column,$
$\}$
```
A `$\ColumnProof_c$` is an `$\OctTreeProof_c$` adjoined with an additional field `$\column_c$`, the Merkle challenge `$c$`'s label in each Stacked-DRG layer. A valid `$\ColumnProof$` has `$\ColumnProof\dot\leaf: \ColumnDigest = \Poseidon_{11}(\ColumnProof\dot\column) \thin$`.

`$\ParentLabels_{\mathbf{u}_\drg}: \Label^{[d_\drg]} = [\Label_{u_\drg, l_0} \mid \forall u_\drg \in \mathbf{u}_\drg]$`\
`$\ParentLabels_{\mathbf{u}_\total}: \Label^{[d_\total]} = [\Label_{u_\drg, l} \mid \forall u_\drg \in \mathbf{u}_\drg] \concat [\Label_{u_\exp, l - 1} \mid \forall u_\exp \in \mathbf{u}_\exp]$`\
The arrays of a node `$v$`'s (where `$v$` is in Stacked-DRG layer `$l_v$`) DRG and total parent labels respectively (DRG parent labels are in layer `$l_v$`, expander parent labels are in layer `$l_v - 1$`). `$\ParentLabels_{\mathbf{u}_\drg}$` is only called while labeling nodes in layer `$l_0 = 0 \thin$`.

`$\ParentLabels_{\mathbf{u}_\drg}^\star: \Label^{[N_\parentlabels]} = \ParentLabels_{\mathbf{u}_\drg}\dot\repeattolength(N_\parentlabels)$`\
`$\ParentLabels_{\mathbf{u}_\total}^\star: \Label^{[N_\parentlabels]} = \ParentLabels_{\mathbf{u}_\total}\dot\repeattolength(N_\parentlabels)$`\
The superscript `$^\star$` denotes that `$\ParentLabels$` has been expanded to length `$N_\parentlabels \thin$`.

`$\R_{\porepchallenges, \batch}: \Byte^{[32]}$`\
Randomness used to generate the challenge set for a batch of PoRep proofs.

`$\PorepChallenges_{R, k}: (\NodeIndex \setminus 0)^{[N_{\porepchallenges / k}]}$`\
The set of PoRep challenges for a replica `$R$`'s partition-`$k$` partition proof. The first node index `$0$` is not challenged in PoRep's (the operator `$\setminus$` is set subtraction).

```text
$\struct \PorepChallengeProof_c\ \{$
$\quad \TreeDProof_c,$
$\quad \ColumnProof_c,$
$\quad \TreeRProof_c,$
$\quad \ParentColumnProofs_{\mathbf{u}_\total}: {\ColumnProof_u}^{[d_\total]},$
$\}$
```
The proof generated for each Merkle challenge `$c$` in a PoRep partition proof. The field `$\ParentColumnProofs_{\mathbf{u}_\total}$` stores the `$\ColumnProof_u$` for each parent `$u \in \mathbf{u}_\total$` of the challenge `$c$`.

`$\PorepPartitionProof_{R, k}: \PorepChallengeProof^{[N_{\porepchallenges / k}]}$`\
A single PoRep partition proof for partition `$k$` in a replica `$R$`'s PoRep proof batch.

### PoSt

`$\R_{\postchallenges, \batch \thin \aww}: \Fq$`\
Randomness used to generate the challenge set for a batch of PoSt proofs.

```text
$\struct \PostReplica_\P\ \{$
$\quad \TreeR,$
$\quad \CommC,$
$\quad \CommCR,$
$\quad \SectorID,$
$\}$
```
A replica `$R$` that a PoSt prover has access to. All fields are associated with `$R$` at the time of PoSt proof batch generation.

`$\PostReplicas_{\P, \batch \thin \aww}: {\PostReplica}^{[*]}$`\
The set of all distinct replicas that a PoSt prover `$\P$` has knowledge of at time of a Winning or Window batch proof generation. PoSt provers in the Filecoin network may have different sized replica sets, thus `$\PostReplicas$` is arbitrarily sized `$*$`.

```text
$\PostReplicas_{\P, k \thin \aww}: {\PostReplica_{\P, \batch}}^{[0 < \ell \leq N_{\postreplicas/k \thin \aww}]}$
$\PostReplicas_{\P, k \thin \aww} =$
$\quad \PostReplicas_{\P, \batch \thin \aww}[k * N_{\postreplicas / k \thin \aww} \ldotdot (k + 1) * N_{\postreplicas / k \thin \aww}]$
```
The `$k^{th}$` distinct slice of a PoSt prover `$\P$`'s total replica set `$\PostReplicas_{\P, \batch}$` used to generate prover's partition-`$k$` Winning or Window PoSt proof in a batch. This set contains all replicas that are challenged in PoSt partition `$k$`. `$\PostReplicas_{\P, k}$` does not contain padded replica proofs. The length of a PoSt prover's total replica set may not be divisible by the number of challenged replica's `$N_{\postreplicas / k \thin \aww}$`, thus the length of `$\PostReplicas_{\P, k}$` is in the range `$[1, N_{\postreplicas / k \thin \aww}] \thin$`.

`$\PostPartitionProof_{k \thin \aww}: {\PostReplicaProof_{R \thin \aww}}^{[N_{\postreplicas/k \thin \aww}]}$`\
A PoSt partition proof generated by a PoSt prover for their `$k^{th}$` partition proof in their current batch of Winning or Window PoSt proofs. Each `$\PostReplicaProof$` in the partition proof is associated with a unique challenged replica `$R$` (unique across the entire batch). A `$\PostPartitionProof$` may contain padded replica proofs to ensure that the partition proof has length `$N_{\postreplicas / k \thin \aww} \thin$`.

```text
$\struct \PostReplicaProof_{R \thin \aww}\ \{$
$\quad \TreeRProofs: \TreeRProof^{[N_{\challenges/R \thin \aww}]} \thin,$
$\quad \CommC,$
$\}$
```
The proof for single replica `$R$` challenged in a PoSt partition proof. All fields are associated with `$R$`.

```text
$\struct \PostReplica_\V\ \{$
$\quad \SectorID,$
$\quad \CommCR,$
$\}$
```
The public information known to a PoSt verifier `$\V$` for each challenged replica `$R$` (distinct up to the PoSt batch) of the PoSt prover. `$\SectorID$` and `$\CommCR$` are associated with the replica `$R$`.

`$\PostReplicas_{\V, k \thin \aww}: {\PostReplica_\V}^{[0 < \ell \leq N_{\postreplicas / k \thin \aww}]}$`\
The public information known to PoSt verifier `$\V$` for each distinct replica `$R$` in a PoSt prover's partition-`$k$` replica set `$\Replicas_{\P, k} \thin$`. The length of the partition's replica set is the number of unique replica's used to generate `$\PartitionProof_{P, k \thin \aww}$` which may be less than `$N_{\postreplicas / k \thin \aww} \thin$`.

## Type Conversions

`$x \as \mathbb{T}$`\
Converts a value `$x$` of type `$\mathbb{X}$` into type `$\mathbb{T}$`. The `$\as$` keyword **takes precedence over arithmetic and bitwise operations**, for example `$a * b \as \mathbb{T}$` casts `$b$` to type `$\mathbb{T}$` before to multiplying by `$a$`.

`$\Safe \as \Fq$`\
`$\Safe \as \Fqsafe$`\
Converts a 32-byte array where only the lowest 254 bits are utilized into a prime field element. A `$\Safe$` byte array is guaranteed to represent a valid field element.

`$\Fq \as \Byte^{[32]}$`\
Converts a 32-byte value to a prime field element. This conversion is safe as all field elements can be represented using `$\ell_\Fq^\bit = 255$` bits.

`$\beencode(\mathbb{T}) \rightarrow \Byte^{[n]}$`\
`$\leencode(\mathbb{T}) \rightarrow \Byte^{[n]}$`\
Big and little-endian encoding (big and little endian with respect to the byte order of the produced byte array) of a value having type `$\mathbb{T}$` into an array of `$n$` bytes. 32-bit integers `$\u{32}$`'s are encoded into 4 bytes and 64-bit integers `$\u{64}$`'s are encoded into 8 bytes.

`$\bedecode(\Byte^{[n]}) \rightarrow \mathbb{T}$`\
`$\ledecode(\Byte^{[n]}) \rightarrow \mathbb{T}$`\
Decodes big and little-endian byte arrays into a value of type `$\mathbb{T}$`.

`$\Bit_\Le^{[0 < n \leq \ell_\Fq^\bit]} \quad \as \Fq$`\
`$\Bit_\Le^{[0 < n \leq \ell_\Fqsafe^\bit]} \thin \as \Fqsafe$`\
Casts a little-endian array of `$n$` bits into a field element (`$\ell_\Fq^\bit = 255$`) and safe field element (`$\ell_\Fqsafe^\bit = 254$`) respectively. Equivalent to writing:

```text
$\bi \line{1}{\bi}{\bits_\Le: \Bit^{[n]} = \text{<some bits>}}$
$\bi \line{2}{\bi}{\fieldelement: \Fq = 0}$
$\bi \line{3}{\bi}{\for i \in [n]:}$
$\bi \line{4}{\bi}{\quad \fieldelement \mathrel{+}= 2^i * \bits_\Le[i]}$
```

`$\bits_{[\auxb, \Le]} \as \Fq$`\
Creates a field element from an array of allocated bits in little-endian order. Equivalent to writing:

```text
$\bi \line{1}{\bi}{\bits_{[\auxb, \Le]}: \CircuitBit^{[0 < n \leq \ell_\Fq^\bit]} = \text{<some allocated bits>}}$
$\bi \line{2}{\bi}{\fieldelement: \Fq = 0}$
$\bi \line{3}{\bi}{\for i \in [n]:}$
$\bi \line{4}{\bi}{\quad \fieldelement \mathrel{+}= 2^i * \bits_{[\auxb, \Le]}[i]\dot\int}$
```

## R1CS and Circuit Notation and Types

`$\RCS$`\
The type used to represent an instance of a rank-1 constraint system. Each `$\RCS$` instance can be thought of as a structure containing two vectors, the primary and auxiliary assignments, and a system of quadratic constraints (a vector where each element is an equality of the form `$\LinearCombination * \LinearCombination = \LinearCombination$`, and each constraint polynomial's variables are values allocated within an `$\RCS$` assignments vectors). The `$\RCS$` type is left opaque as its implementation is beyond the scope of this document.

`$\RCS\dot\one_\pubb = \CircuitVal\ \{ \index_\text{primary}: 0, \int: 1\ \}$`\
Every `$\RCS$` instance is instantiated with its first primary assignment being the multiplicative identity.

```text
$\RCS\dot\publicinput(\Fq)$
$\RCS\dot\privateinput(\Fq)$
$\RCS\dot\assert(\LinearCombination * \LinearCombination = \LinearCombination)$
```
The `$\RCS$` types has three methods, one for adding primary assignments, one for adding auxiliary assignments, and one for adding constraints.

```text
$\struct \CircuitVal\ \{$
$\quad \index_{\langle \text{primary|auxiliary} \rangle}: \mathbb{N},$
$\quad \int: \Fq,$
$\}$
```
An instance of `$\CircuitVal$` is a reference to a single field element, of value `$\int$`, allocated within a constraint system. The field `$\int$` is a copy of the allocated value, arithmetic on `$\int$` is not circuit arithmetic. The field `$\index$` refers to the index of an allocated integer (a copy of `$\int$`) in either the primary or auxiliary assignments vectors of an `$\RCS$` instance. Every unique wire in a circuit has a unique (up to the instance of `$\RCS$`) `$\CircuitVal\dot\index$` in either the primary or auxiliary assignments vectors.

```text
$\struct \CircuitVal_\Bit \equiv \CircuitVal\ \{$
$\quad \index_{\langle \text{primary|auxiliary} \rangle}: \mathbb{N},$
$\quad \int: \Fq \in \Bit,$
$\}$
```
A reference to an allocated bit, a boolean constrained value, within an `$\RCS$` instance.

`$\CircuitBitOrConst \equiv \{ \CircuitBit,\thin \Bit \}$`\
The type `$\CircuitBitOrConst$` is an enum representing either an allocated bit or a constant unalloacted bit.

`$\deq$`\
The "diamond-equals" sign shows that a value as been allocated within `$\RCS$`. If an assignment `$=$` operates on `$\CircuitVal{\bf \sf s}$` and does not have a diamond, then no value was allocated in a circuit.

`$\value_\auxb: \CircuitVal \deq \RCS\dot\privateinput(\value \as \Fq)$`\
The subscript `$_\auxb$` denotes `$\value$` as being allocated within `$\RCS$` and located in the auxiliary assignments vector.

`$\value_\pubb: \CircuitVal \deq \RCS\dot\publicinput(\value \as \Fq)$`\
The subscript `$_\pubb$` denotes an allocated value as being in an `$\RCS$` primary assignments vector. The function `$\RCS\dot\publicinput$` adds a public input `$\value$` to the primary assignments vector (denoted `$\value_\pubb$`), allocates `$\value$` within the auxiliary assignments vector (denoted `$\value_\auxb$`), and adds an equality constraint checking that the SNARK prover's auxiliary assignment `$\value_\auxb$` is equal to the verifiers's public input `$\value_\pubb$`.

`$\value_\aap: \CircuitVal$`\
The allocated `$\value$` may be either an auxiliary or primary assignment.

`$\bit_\pubb: \CircuitBit \deq \RCS\dot\publicinput(\Bit \as \Fq)$`\
`$\bit_\auxb: \CircuitBit \deq \RCS\dot\privateinput(\Bit \as \Fq)$`\
The subscript `$_\Bit$` in `$\CircuitBit$` denotes an allocated value that has been boolean constrained. The functions `$\RCS\dot\publicinput$` and `$\RCS\dot\privateinput$` add the boolean constraint `$(1 - \bit_\aap) * (\bit_\aap) = 0$` opaquely based upon whether or not the return type is `$\CircuitVal$` or `$\CircuitBit \thin$`.

`$\bits_{[\auxb]}: \CircuitBit^{[*]}$`\
The subscript `$_{[\auxb]}$` denotes a value as being an array of allocated bits (boolean constrained circuit values).

`$\bits_{[\auxb + \constb]}: \CircuitBitOrConst^{[*]}$`\
The subscript `$_{[\auxb + \constb]}$` denotes a value as being an array of allocated bits and constant/unallocated boolean values `$\Bit$`.

`$\bits_{[\auxb, \Le]}: \CircuitBit^{[*]}$`\
The subscript `$_{[\auxb, \Le]}$` denotes an array of allocated bits as having little-endian bit order (least significant bit is at index 0 and the most significant bit is the last element in the array).

`$\bits_{[\auxb, \lebytes]}: \CircuitBit^{[*]}$`\
The subscript `$_{[\auxb, \lebytes]}$` denotes an array of allocated bits as having little-endian byte order (least significant byte first where the most significant bit of each byte is at index 0 in that byte's 8-bit slice of the array).

`$\arr_{[[\auxb]]}: {\CircuitBit^{[m]}}^{[n]}$`\
The subscript `$_{[[\auxb]]}$` denotes an array where each element is an array of allocated bits.

`$\LinearCombination(\CircuitVal_0, \ldots, \CircuitVal_n)$`\
Represents an unevaluated linear combination of allocated value variables and unallocated constants, for example the degree-1 polynomial `$2 * \CircuitVal_0 + \CircuitVal_1 + 5$` (where `$2$` and `$5$` are unallocated constants) is a `$\LinearCombination \thin$`.

`$\lc: \LinearCombination \equiv 0$`\
`$\lc: \LinearCombination \equiv \CircuitVal_0 + \CircuitVal_1$`\
`$\lc: \LinearCombination \equiv \sum_{i \in [3]}\ 2^i * \CircuitVal$`\
When linear combinations are instantiated they are not evaluated to a single integer value and stored in `$\lc$`. The equivalency notation `$\equiv$` is used show that a linear combination is not evaluated to a single value, but represents a symbolic polynomial over allocated `$\CircuitVal$`'s. In the above examples the values `$0$` and `$2^i$` are unallocated constants.

## Protocol Assumptions

Values are chosen for Protocol constants such that the following are true:

* The sector size `$\ell_\sector^\byte$` is divisible by the node size `$\ell_\node^\byte$`.
* The number of nodes always fits an unsigned 32-bit integer `$N_\nodes \leq 2^{32}$`, and thus node indexes are representable using 32-bit unsigned integers.
* Every distinct contiguous 32-byte slice of a sector `$D_i \in D$` represents a valid prime field element `$\Fq$` (every `$D_i$` is a valid `$\Safe$`, this is accomplished via sector preprocessing).

## Hash Functions

`$\Sha{256}(\Bit^{[*]}) \rightarrow \Byte^{[32]}$`\
The `$\textsf{Sha256}$` hash function operates on preimages of an arbitrary number of bits and returns 32-byte array (256 bits).

`$\Sha{254}(\Bit^{[*]}) \rightarrow \Safe$`\
`$\Sha{254}_2({\Byte^{[32]}}^{[2]}) \rightarrow \Safe$`\
The `$\Sha{254}$` hash functions are identical to `$\Sha{256}$` except that the last two bits of the `$\Sha{256}$` 256-bit digest are zeroed out:

```text
$\bi \Sha{254}(x) \equiv \Sha{256}(x)[\ldotdot 31] \concat (\Sha{256}(x)[31] \AND 1^{[6]})$
```

The hash function `$\Sha{254}$` operates on preimages of unsized bit arrays while `$\Sha{254}_2$` operates on preimages of two 32-byte values.

`$\Poseidon(\Bit^{[*]}) \rightarrow \Fq$`\
`$\Poseidon_2(\Fq^{[2]}) \rightarrow \Fq$`\
`$\Poseidon_8(\Fq^{[8]}) \rightarrow \Fq$`\
`$\Poseidon_{11}(\Fq^{[11]}) \rightarrow \Fq$`\
The `$\Poseidon$` hash functions operate on preimages of an arbitrary number of bits as well as preimages containing 2, 8, and 11 field elements.

## Naming Differences

This specification document deviates in one major way from the naming conventions used within the [Filecoin proofs implementation](https://github.com/filecoin-project/rust-fil-proofs/tree/master/storage-proofs). That is what the code calls [`tree_r_last`](https://github.com/filecoin-project/rust-fil-proofs/blob/b288702362e8f14ee0a68fb030774f339266e9a6/storage-proofs/porep/src/stacked/vanilla/proof.rs#L690) and [`comm_r_last`](https://github.com/filecoin-project/rust-fil-proofs/blob/b288702362e8f14ee0a68fb030774f339266e9a6/storage-proofs/porep/src/stacked/vanilla/proof.rs#L1083), this document calls `$\TreeR$` and `$\CommR$` respectively. What the code calls [`comm_r`](https://github.com/filecoin-project/rust-fil-proofs/blob/b288702362e8f14ee0a68fb030774f339266e9a6/storage-proofs/porep/src/stacked/vanilla/proof.rs#L1084), this document calls `$\CommCR$`.

In the code, `tree_r_last` is built over the replica, this is why this document uses the notation `$\TreeR$`.

In the code, `comm_r` is not the root of the Merkle tree built over the replica `tree_r_last`, but is the hash of `tree_c.root` and `tree_r_last.root`. This is why this specification document has changed `comm_r`'s name to `$\CommCR$`.

`tree_r_last` `$\mapsto \TreeR$`\
`comm_r_last` `$\mapsto \CommR$`\
`comm_r` `$\mapsto \CommCR$`

## Bit Ordering

### Little-Endian Bits

A bit array whose least significant bit is at index `$0$` and where each subsequent index points to the next most significant bit has little-endian bit order. Little endian bit arrays are denoted by the subscript `$_\Le\thin$`, for example:

```text
$\bi \bits_\Le: \Bit^{[n]} = \lebinrep{x}$
```

An unsigned integer `$\int$`'s little-endian `$n$`-bit  binary representation `$\lebinrep{\int}$`, where `$n = \lceil \log_2(\int) \rceil$`, is defined as:

```text
$\bi \lebinrep{\int}: \Bit^{[n]} = [(\int \gg i) \AND 1 \mid \forall i \in [n]]$
```

An unsigned integer `$\int$`'s little-endian binary representation can be padded with 0's at the most significant end to yield an `$n'$`-bit binary representation, where `$n' > n \thin$`:

```text
$\bi \lebinrep{\int}: \Bit^{[n']} = \lebinrep{\int} \concat 0^{[n' - n]}$
```

An unsigned integer `$\int$`'s little-endian  bit string is the reverse of its big-endian bit string:

```text
$\bi \lebinrep{\int} = \reverse(\bebinrep{\int})$
$\bi \reverse(\lebinrep{\int}) = \bebinrep{\int}$
```

### Big-Endian Bits

A bit array whose most significant bit is at index `$0$` and where each subsequent index points to the next least significant bit has big-endian bit order. Big-endian bit strings are denotes by the subscript `$_\be \thin$`, for example:

```text
$\bi \bits_\be: \Bit^{[n]} = \bebinrep{x}$
```

An unsigned integer `$\int$`'s `$n$`-bit big-endian bit string, where `$n = \lceil \log_2(\int) \rceil$`, is defined as:

```text
$\bi \bebinrep{\int}: \Bit^{[n]} = [(\int \gg (n - i - 1)) \AND 1 \mid \forall i \in [n]]$
```

An unsigned integer's big-endian bit string is the reverse of its little-endian bit string:

```text
$\bi \bebinrep{\int} \equiv \reverse(\lebinrep{\int})$
$\bi \reverse(\bebinrep{\int}) = \lebinrep{\int}$
```

### Little-Endian Bytes

The `$\lebytes$` bit order signifies that an `$n$`-bit bit string (where `$n$` is evenly divisible by 8) has bit ordering such that: each distinct 8-bit slice of the array is more significant the previous (the first byte is the least significant) and each byte has big-endian bit order (the first bit in each byte is the most significant with respect to that 8-bit slice). Bit strings of `$\lebytes$` bit order have length `$n = \lceil \log_2(\int) \rceil$` where `$n$` is a multiple of 8. Bit strings having `$\lebytes$` bit order are denoted using the subscript `$_\lebytes \thin$`.

An unsigned integer `$\int$` represented in binary using `$\lebytes$` bit order is defined as:

```text
$\bi \lebytesbinrep{\int} \bi : \Bit^{[n]} = \byte_0 \concat \byte_1 \concat \ldots \concat \byte_{n / 8}$
$\bi \byte_0 \quad\quad\bi\thin\thin\thin : \Bit^{[8]} = [(\LSByte_0, \MSBit), \dots, (\LSByte_0, \LSBit)]$
$\bi \byte_1 \quad\quad\bi\thin\thin\thin : \Bit^{[8]} = [(\LSByte_1, \MSBit), \dots, (\LSByte_1, \LSBit)]$
$\bi \dots$
$\bi \byte_{n / 8} \quad\quad\thin\thin : \Bit^{[8]} = [(\MSByte, \MSBit), \dots, (\MSByte, \LSBit)]$
```

where the integer's least significant byte `$\LSByte_0$` is the first byte `$\byte_0$` in the `$\lebytes$` bit string, the integer's second least significant byte `$\LSByte_1$` is the second byte in the bit string `$\byte_0$`, and so one until the last byte in the bistring is the most significant byte `$\MSByte$` with respect to the intger's value. Each byte in the bit string has big-endian bit order, each bytes most significant bit `$\MSBit$` is first and least-significant bit `$\LSBit$` is last. `$n / 8$` is the number of distinct 8-bit slices of `$\int$`'s `$n$`-bit binary representation.

The `$\lebytes$` bit order is used to represent a field element `$\Fq$` as a 256-bit `$\Sha{256}$` input block. Because `$\Fq$` and `$\Fqsafe$` have bit lengths of 255 and 254 respectively, one or two zero bits are padded onto the most significant end of `$\lebinrep{\Fq}$` and `$\lebinrep{\Fqsafe}$` to fill the 256-bit SHA block. The padding operation occurs when the bit string has little-endian bit order `$\bits_\Le \concat 0^{[\len(\bits_\Le) \thin\MOD\thin 8]}$` before it is converted to `$\lebytes$` bit order.

### Little-Endian Bits to Little-Endian Bytes

**Implementation:** [`storage_proofs::core::util::reverse_bit_numbering()`](https://github.com/filecoin-project/rust-fil-proofs/blob/3cc8b5acb742d1469e10949fadc389806fa19c8c/storage-proofs/core/src/util.rs#L130)

```text
$\overline{\underline{\Function \lebitstolebytes(\bits_\Le: \Bit^{[m]}) \rightarrow \Bit^{[n]}}}$
$\line{1}{\bi}{\bits_\Le: \Bit^{[n]} = \bits_\Le \concat 0^{[m \thin\MOD\thin 8]}}$
$\line{2}{\bi}{\bits_\lebytes: \Bit^{[n]} = [\ ]}$
$\line{3}{\bi}{\for i \in [n / 8]:}$
$\line{4}{\bi}{\quad \byte_\Le: \Bit^{[8]} = \bits_\Le[i * 8 \thin\ldotdot\thin (i + 1) * 8]}$
$\line{5}{\bi}{\quad \byte_\be: \Bit^{[8]} = \reverse(\byte_\Le)}$
$\line{6}{\bi}{\quad \bits_\lebytes\dot\extend(\byte_\be)}$
$\line{7}{\bi}{\return \bits_\lebytes}$
```

### Little-Endian Bytes to Little-Endian Bits

The length `$n$` of `$\bits$` must be divisible by 8 (must contain an integer number of bytes).

```text
$\overline{\underline{\Function \lebytestolebits(\bits_\lebytes: \Bit^{[n]}) \rightarrow \Bit^{[n]}}}$
$\line{1}{\bi}{\assert(n \MOD 8 = 0)}$
$\line{2}{\bi}{\bits_\Le: \Bit^{[n]} = [\ ]}$
$\line{3}{\bi}{\for i \in [n / 8]:}$
$\line{4}{\bi}{\quad \byte_\be: \Bit^{[8]} = \bits_\lebytes[i * 8 \thin\ldotdot\thin (i + 1) * 8]}$
$\line{5}{\bi}{\quad \byte_\Le: \Bit^{[8]} = \reverse(\byte_\be)}$
$\line{6}{\bi}{\quad \bits_\Le\dot\extend(\byte_\Le)}$
$\line{7}{\bi}{\return \bits_\Le}$
```

## Filecoin Proofs Terminology

### PoRep and PoSt

The Filecoin protocol has three proof variants: PoRep (proof-of-replication), Winning PoSt (proof-of-spacetime), and Window PoSt. The two PoSt variants, Winning and Window, are identical aside from their number of Merkle challenges and number of replica `$\TreeR$`'s for which the Merkle challenges are made.

A PoRep proof (and PoRep proof batch) are made for a single replica, whereas a PoSt proof (and proof batch) are made using one or more replicas determined by the constants `$N_{\replicas / k \thin \aww} \thin$`.

### Vanilla Proofs v.s. SNARKs

Each proof variant is proven in two ways: vanilla (non-SNARK) and SNARK.

Vanilla proofs are used to instantiate an instance of a proof variant's circuit. A circuit instance is then used to generate a `$\GrothProof$`.

### Partitions Proofs and Batches

Each proof variant's vanilla proof is called a partition proof, `$\PorepPartitionProof$` and `$\PostPartitionProof$`. A PoRep or PoSt prover generates multiple partition proofs simultaneously, called a vanilla proof batch. The number of partition proofs per proof variant batch are given by the constants: `$N_{\poreppartitions / \batch}$`, `$N_{\postpartitions / \batch, \P, \winning}$`, and `$N_{\postpartitions / \batch, \P, \window} \thin$`. A SNARK is generated for each vanilla proof in a vanilla proof batch resulting in a batch of corresponding SNARK proofs.

The function `$\createporepbatch_R$` shows how a PoRep batch proof can be made for a replica `$R$`, where `$R$` (and its associated trees, commitments, and labeling) were output by the replication process. A similar process is used to produce Winning and Window PoSt batches.

```text
$\overline{\Function \createporepbatch_R(\qquad\qquad\qquad\qquad\quad}$
$\quad \ReplicaID,$
$\quad \TreeD,$
$\quad \TreeC,$
$\quad \TreeR,$
$\quad \CommD,$
$\quad \CommC,$
$\quad \CommR,$
$\quad \CommCR,$
$\quad \Labels,$
$\quad \R_{\porepchallenges, \batch} \thin,$
$\underline{) \rightarrow (\PorepPartitionProof, \GrothProof)^{[N_{\poreppartitions / \batch}]}}$
$\line{1}{\bi}{\batch: (\PorepPartitionProof, \GrothProof)^{[N_{\poreppartitions / \batch}]} = [\ ]}$
$\line{2}{\bi}{\for k \in [N_{\poreppartitions / \batch}]:}$
$\line{3}{\bi}{\quad \PorepPartitionProof_k = \createvanillaporepproof(}$
$\quad\quad\quad k,$
$\quad\quad\quad \ReplicaID,$
$\quad\quad\quad \TreeD,$
$\quad\quad\quad \TreeC,$
$\quad\quad\quad \TreeR,$
$\quad\quad\quad \Labels,$
$\quad\quad\quad \R_{\porepchallenges, \batch} \thin,$
$\quad\quad\thin )$
$\line{4}{\bi}{\quad \RCS_k = \createporepcircuit(}$
$\quad\quad\quad \PorepPartitionProof_k,$
$\quad\quad\quad k,$
$\quad\quad\quad \ReplicaID,$
$\quad\quad\quad \CommD,$
$\quad\quad\quad \CommC,$
$\quad\quad\quad \CommR,$
$\quad\quad\quad \CommCR,$
$\quad\quad\quad \R_{\porepchallenges, \batch} \thin,$
$\quad\quad\thin )$
$\line{5}{\bi}{\quad \GrothProof_k = \creategrothproof(\GrothEvaluationKey_\porep, \RCS_k)}$
$\line{6}{\bi}{\quad \batch\dot\push((\PorepPartitionProof_k, \GrothProof_k))}$
$\line{7}{\bi}{\return \batch}$
```
