------------------------ MODULE PieceInclusionProof ------------------------
EXTENDS Integers, TLC, Sequences, FiniteSets

SIZE == 256 \* Must be a power of 2. TODO: Make CONSTANT and add ASSUME.
HEIGHT == 9 \* TODO: Calculate this from SIZE, as log2(SIZE) + 1.

(*--algorithm PieceInclusionProof

variables
    HashCounter = -1;
    HashRecord = <<>>;

macro RepCompress(a, b, height, var) begin
    if (<<a, b, height>> \in DOMAIN HashRecord) then
        var := HashRecord[<<a, b, height>>];
    else
        HashCounter := HashCounter - 1;
        HashRecord := (<<a, b, height>> :> HashCounter) @@ HashRecord;    
        var := HashCounter
    end if;
end macro;

process test_hash = "test hash"
variables h1, h2, h3;

begin
    L1:
        RepCompress(1, 2, 0, h1);
    L2:        
        RepCompress(1, 2, 0, h2);
    L3:
        RepCompress(2, 1, 0, h3);
        
        assert h1 = h2;
        assert h1 /= h3;
end process;

fair process merkle_tree = "merkle tree"
variables h,
    input, row, rowSize, nextRow, index, proof_element, root, challenge,
    cursor_index, cursor_row, cursor_element, proof_index, challenge_path_acc, place_acc,
    rows = <<>>;
    height = -1;
    
    proof_path = <<>>;
    proof_elements = <<>>;
begin
    BuildTree:
        input := [i \in 1..SIZE |-> i];
        row := input;
        rows := <<>>;
    RowLoop:
        height := height + 1;
        rows := Append(rows, row);
        \* It would be nice to make this assert an invariant, but how do we make an invariant
        \* over a process variable?
        assert height > 1 => Len(rows[height-1]) = 2 * Len(rows[height]);
        
        nextRow := <<>>;
        index := 1;
        rowSize := Cardinality(DOMAIN row);    
        if rowSize > 1 then
            HashRow:
                RepCompress(row[index], row[index+1], height, h);
                nextRow := Append(nextRow, h);
            Advance: index := index + 2;
            if index < Cardinality(DOMAIN row) then
                goto HashRow;
            else
                row := nextRow;
            end if;
            Repeat: goto RowLoop;
        else
            assert Len(rows) = HEIGHT;
        end if;
    Proofs:
        challenge := 1;
    MakeProof:
        cursor_index := challenge;
        cursor_row := 1;
        cursor_element := rows[cursor_row][cursor_index];
        proof_path := <<>>;
        proof_elements := <<>>;
    S1:
        if cursor_index % 2 = 1 then
            proof_path := Append(proof_path, FALSE);
            proof_element := rows[cursor_row][cursor_index+1];
            RepCompress(rows[cursor_row][cursor_index],
                        proof_element,
                        cursor_row - 1,
                        cursor_element);
        else
            proof_path := Append(proof_path, TRUE);
            proof_element := rows[cursor_row][cursor_index-1];
            RepCompress(proof_element,
                        rows[cursor_row][cursor_index],
                        cursor_row - 1,
                        cursor_element);
        end if;
        
        proof_elements := Append(proof_elements, proof_element);
    ProofLoop:
        cursor_row := cursor_row + 1;
        cursor_index := (cursor_index + 1) \div 2;
        if cursor_row < Len(rows) then
            goto S1;
        end if;
    FinishProof:
        root := rows[Len(rows)][1];
        
    CheckProof:
        proof_index := 1;
        height := 0;
        cursor_index := challenge;
        cursor_element := rows[height+1][cursor_index];
        challenge_path_acc := 0;
        place_acc := 1;
    ProofCheckLoop:
        if proof_path[proof_index] then
            RepCompress(proof_elements[proof_index],
                        cursor_element,
                        height,
                        cursor_element);
            challenge_path_acc := challenge_path_acc + place_acc;
        else
           RepCompress(cursor_element,
                       proof_elements[proof_index],
                       height,
                       cursor_element); 
        end if;
        place_acc := place_acc * 2;
        
        proof_index := proof_index + 1;
        height := height + 1;
        if height < Len(proof_elements) then
            goto ProofCheckLoop;
        end if;
    
    CheckRoot:
        assert cursor_element = root;
        assert challenge_path_acc = challenge - 1; \* challenges are 1-indexed because TLA+.
    
    IncrementChallenge:
        challenge := challenge + 1;
        if challenge <= Len(input) then
            goto MakeProof;
        end if;
end process;

end algorithm; *)

\* BEGIN TRANSLATION
CONSTANT defaultInitValue
VARIABLES HashCounter, HashRecord, pc, h1, h2, h3, h, input, row, rowSize, 
          nextRow, index, proof_element, root, challenge, cursor_index, 
          cursor_row, cursor_element, proof_index, challenge_path_acc, 
          place_acc, rows, height, proof_path, proof_elements

vars == << HashCounter, HashRecord, pc, h1, h2, h3, h, input, row, rowSize, 
           nextRow, index, proof_element, root, challenge, cursor_index, 
           cursor_row, cursor_element, proof_index, challenge_path_acc, 
           place_acc, rows, height, proof_path, proof_elements >>

ProcSet == {"test hash"} \cup {"merkle tree"}

Init == (* Global variables *)
        /\ HashCounter = -1
        /\ HashRecord = <<>>
        (* Process test_hash *)
        /\ h1 = defaultInitValue
        /\ h2 = defaultInitValue
        /\ h3 = defaultInitValue
        (* Process merkle_tree *)
        /\ h = defaultInitValue
        /\ input = defaultInitValue
        /\ row = defaultInitValue
        /\ rowSize = defaultInitValue
        /\ nextRow = defaultInitValue
        /\ index = defaultInitValue
        /\ proof_element = defaultInitValue
        /\ root = defaultInitValue
        /\ challenge = defaultInitValue
        /\ cursor_index = defaultInitValue
        /\ cursor_row = defaultInitValue
        /\ cursor_element = defaultInitValue
        /\ proof_index = defaultInitValue
        /\ challenge_path_acc = defaultInitValue
        /\ place_acc = defaultInitValue
        /\ rows = <<>>
        /\ height = -1
        /\ proof_path = <<>>
        /\ proof_elements = <<>>
        /\ pc = [self \in ProcSet |-> CASE self = "test hash" -> "L1"
                                        [] self = "merkle tree" -> "BuildTree"]

L1 == /\ pc["test hash"] = "L1"
      /\ IF (<<1, 2, 0>> \in DOMAIN HashRecord)
            THEN /\ h1' = HashRecord[<<1, 2, 0>>]
                 /\ UNCHANGED << HashCounter, HashRecord >>
            ELSE /\ HashCounter' = HashCounter - 1
                 /\ HashRecord' = (<<1, 2, 0>> :> HashCounter') @@ HashRecord
                 /\ h1' = HashCounter'
      /\ pc' = [pc EXCEPT !["test hash"] = "L2"]
      /\ UNCHANGED << h2, h3, h, input, row, rowSize, nextRow, index, 
                      proof_element, root, challenge, cursor_index, cursor_row, 
                      cursor_element, proof_index, challenge_path_acc, 
                      place_acc, rows, height, proof_path, proof_elements >>

L2 == /\ pc["test hash"] = "L2"
      /\ IF (<<1, 2, 0>> \in DOMAIN HashRecord)
            THEN /\ h2' = HashRecord[<<1, 2, 0>>]
                 /\ UNCHANGED << HashCounter, HashRecord >>
            ELSE /\ HashCounter' = HashCounter - 1
                 /\ HashRecord' = (<<1, 2, 0>> :> HashCounter') @@ HashRecord
                 /\ h2' = HashCounter'
      /\ pc' = [pc EXCEPT !["test hash"] = "L3"]
      /\ UNCHANGED << h1, h3, h, input, row, rowSize, nextRow, index, 
                      proof_element, root, challenge, cursor_index, cursor_row, 
                      cursor_element, proof_index, challenge_path_acc, 
                      place_acc, rows, height, proof_path, proof_elements >>

L3 == /\ pc["test hash"] = "L3"
      /\ IF (<<2, 1, 0>> \in DOMAIN HashRecord)
            THEN /\ h3' = HashRecord[<<2, 1, 0>>]
                 /\ UNCHANGED << HashCounter, HashRecord >>
            ELSE /\ HashCounter' = HashCounter - 1
                 /\ HashRecord' = (<<2, 1, 0>> :> HashCounter') @@ HashRecord
                 /\ h3' = HashCounter'
      /\ Assert(h1 = h2, "Failure of assertion at line 34, column 9.")
      /\ Assert(h1 /= h3', "Failure of assertion at line 35, column 9.")
      /\ pc' = [pc EXCEPT !["test hash"] = "Done"]
      /\ UNCHANGED << h1, h2, h, input, row, rowSize, nextRow, index, 
                      proof_element, root, challenge, cursor_index, cursor_row, 
                      cursor_element, proof_index, challenge_path_acc, 
                      place_acc, rows, height, proof_path, proof_elements >>

test_hash == L1 \/ L2 \/ L3

BuildTree == /\ pc["merkle tree"] = "BuildTree"
             /\ input' = [i \in 1..SIZE |-> i]
             /\ row' = input'
             /\ rows' = <<>>
             /\ pc' = [pc EXCEPT !["merkle tree"] = "RowLoop"]
             /\ UNCHANGED << HashCounter, HashRecord, h1, h2, h3, h, rowSize, 
                             nextRow, index, proof_element, root, challenge, 
                             cursor_index, cursor_row, cursor_element, 
                             proof_index, challenge_path_acc, place_acc, 
                             height, proof_path, proof_elements >>

RowLoop == /\ pc["merkle tree"] = "RowLoop"
           /\ height' = height + 1
           /\ rows' = Append(rows, row)
           /\ Assert(height' > 1 => Len(rows'[height'-1]) = 2 * Len(rows'[height']), 
                     "Failure of assertion at line 57, column 9.")
           /\ nextRow' = <<>>
           /\ index' = 1
           /\ rowSize' = Cardinality(DOMAIN row)
           /\ IF rowSize' > 1
                 THEN /\ pc' = [pc EXCEPT !["merkle tree"] = "HashRow"]
                 ELSE /\ Assert(Len(rows') = HEIGHT, 
                                "Failure of assertion at line 74, column 13.")
                      /\ pc' = [pc EXCEPT !["merkle tree"] = "Proofs"]
           /\ UNCHANGED << HashCounter, HashRecord, h1, h2, h3, h, input, row, 
                           proof_element, root, challenge, cursor_index, 
                           cursor_row, cursor_element, proof_index, 
                           challenge_path_acc, place_acc, proof_path, 
                           proof_elements >>

HashRow == /\ pc["merkle tree"] = "HashRow"
           /\ IF (<<(row[index]), (row[index+1]), height>> \in DOMAIN HashRecord)
                 THEN /\ h' = HashRecord[<<(row[index]), (row[index+1]), height>>]
                      /\ UNCHANGED << HashCounter, HashRecord >>
                 ELSE /\ HashCounter' = HashCounter - 1
                      /\ HashRecord' = (<<(row[index]), (row[index+1]), height>> :> HashCounter') @@ HashRecord
                      /\ h' = HashCounter'
           /\ nextRow' = Append(nextRow, h')
           /\ pc' = [pc EXCEPT !["merkle tree"] = "Advance"]
           /\ UNCHANGED << h1, h2, h3, input, row, rowSize, index, 
                           proof_element, root, challenge, cursor_index, 
                           cursor_row, cursor_element, proof_index, 
                           challenge_path_acc, place_acc, rows, height, 
                           proof_path, proof_elements >>

Advance == /\ pc["merkle tree"] = "Advance"
           /\ index' = index + 2
           /\ IF index' < Cardinality(DOMAIN row)
                 THEN /\ pc' = [pc EXCEPT !["merkle tree"] = "HashRow"]
                      /\ row' = row
                 ELSE /\ row' = nextRow
                      /\ pc' = [pc EXCEPT !["merkle tree"] = "Repeat"]
           /\ UNCHANGED << HashCounter, HashRecord, h1, h2, h3, h, input, 
                           rowSize, nextRow, proof_element, root, challenge, 
                           cursor_index, cursor_row, cursor_element, 
                           proof_index, challenge_path_acc, place_acc, rows, 
                           height, proof_path, proof_elements >>

Repeat == /\ pc["merkle tree"] = "Repeat"
          /\ pc' = [pc EXCEPT !["merkle tree"] = "RowLoop"]
          /\ UNCHANGED << HashCounter, HashRecord, h1, h2, h3, h, input, row, 
                          rowSize, nextRow, index, proof_element, root, 
                          challenge, cursor_index, cursor_row, cursor_element, 
                          proof_index, challenge_path_acc, place_acc, rows, 
                          height, proof_path, proof_elements >>

Proofs == /\ pc["merkle tree"] = "Proofs"
          /\ challenge' = 1
          /\ pc' = [pc EXCEPT !["merkle tree"] = "MakeProof"]
          /\ UNCHANGED << HashCounter, HashRecord, h1, h2, h3, h, input, row, 
                          rowSize, nextRow, index, proof_element, root, 
                          cursor_index, cursor_row, cursor_element, 
                          proof_index, challenge_path_acc, place_acc, rows, 
                          height, proof_path, proof_elements >>

MakeProof == /\ pc["merkle tree"] = "MakeProof"
             /\ cursor_index' = challenge
             /\ cursor_row' = 1
             /\ cursor_element' = rows[cursor_row'][cursor_index']
             /\ proof_path' = <<>>
             /\ proof_elements' = <<>>
             /\ pc' = [pc EXCEPT !["merkle tree"] = "S1"]
             /\ UNCHANGED << HashCounter, HashRecord, h1, h2, h3, h, input, 
                             row, rowSize, nextRow, index, proof_element, root, 
                             challenge, proof_index, challenge_path_acc, 
                             place_acc, rows, height >>

S1 == /\ pc["merkle tree"] = "S1"
      /\ IF cursor_index % 2 = 1
            THEN /\ proof_path' = Append(proof_path, FALSE)
                 /\ proof_element' = rows[cursor_row][cursor_index+1]
                 /\ IF (<<(rows[cursor_row][cursor_index]), proof_element', (cursor_row - 1)>> \in DOMAIN HashRecord)
                       THEN /\ cursor_element' = HashRecord[<<(rows[cursor_row][cursor_index]), proof_element', (cursor_row - 1)>>]
                            /\ UNCHANGED << HashCounter, HashRecord >>
                       ELSE /\ HashCounter' = HashCounter - 1
                            /\ HashRecord' = (<<(rows[cursor_row][cursor_index]), proof_element', (cursor_row - 1)>> :> HashCounter') @@ HashRecord
                            /\ cursor_element' = HashCounter'
            ELSE /\ proof_path' = Append(proof_path, TRUE)
                 /\ proof_element' = rows[cursor_row][cursor_index-1]
                 /\ IF (<<proof_element', (rows[cursor_row][cursor_index]), (cursor_row - 1)>> \in DOMAIN HashRecord)
                       THEN /\ cursor_element' = HashRecord[<<proof_element', (rows[cursor_row][cursor_index]), (cursor_row - 1)>>]
                            /\ UNCHANGED << HashCounter, HashRecord >>
                       ELSE /\ HashCounter' = HashCounter - 1
                            /\ HashRecord' = (<<proof_element', (rows[cursor_row][cursor_index]), (cursor_row - 1)>> :> HashCounter') @@ HashRecord
                            /\ cursor_element' = HashCounter'
      /\ proof_elements' = Append(proof_elements, proof_element')
      /\ pc' = [pc EXCEPT !["merkle tree"] = "ProofLoop"]
      /\ UNCHANGED << h1, h2, h3, h, input, row, rowSize, nextRow, index, root, 
                      challenge, cursor_index, cursor_row, proof_index, 
                      challenge_path_acc, place_acc, rows, height >>

ProofLoop == /\ pc["merkle tree"] = "ProofLoop"
             /\ cursor_row' = cursor_row + 1
             /\ cursor_index' = ((cursor_index + 1) \div 2)
             /\ IF cursor_row' < Len(rows)
                   THEN /\ pc' = [pc EXCEPT !["merkle tree"] = "S1"]
                   ELSE /\ pc' = [pc EXCEPT !["merkle tree"] = "FinishProof"]
             /\ UNCHANGED << HashCounter, HashRecord, h1, h2, h3, h, input, 
                             row, rowSize, nextRow, index, proof_element, root, 
                             challenge, cursor_element, proof_index, 
                             challenge_path_acc, place_acc, rows, height, 
                             proof_path, proof_elements >>

FinishProof == /\ pc["merkle tree"] = "FinishProof"
               /\ root' = rows[Len(rows)][1]
               /\ pc' = [pc EXCEPT !["merkle tree"] = "CheckProof"]
               /\ UNCHANGED << HashCounter, HashRecord, h1, h2, h3, h, input, 
                               row, rowSize, nextRow, index, proof_element, 
                               challenge, cursor_index, cursor_row, 
                               cursor_element, proof_index, challenge_path_acc, 
                               place_acc, rows, height, proof_path, 
                               proof_elements >>

CheckProof == /\ pc["merkle tree"] = "CheckProof"
              /\ proof_index' = 1
              /\ height' = 0
              /\ cursor_index' = challenge
              /\ cursor_element' = rows[height'+1][cursor_index']
              /\ challenge_path_acc' = 0
              /\ place_acc' = 1
              /\ pc' = [pc EXCEPT !["merkle tree"] = "ProofCheckLoop"]
              /\ UNCHANGED << HashCounter, HashRecord, h1, h2, h3, h, input, 
                              row, rowSize, nextRow, index, proof_element, 
                              root, challenge, cursor_row, rows, proof_path, 
                              proof_elements >>

ProofCheckLoop == /\ pc["merkle tree"] = "ProofCheckLoop"
                  /\ IF proof_path[proof_index]
                        THEN /\ IF (<<(proof_elements[proof_index]), cursor_element, height>> \in DOMAIN HashRecord)
                                   THEN /\ cursor_element' = HashRecord[<<(proof_elements[proof_index]), cursor_element, height>>]
                                        /\ UNCHANGED << HashCounter, 
                                                        HashRecord >>
                                   ELSE /\ HashCounter' = HashCounter - 1
                                        /\ HashRecord' = (<<(proof_elements[proof_index]), cursor_element, height>> :> HashCounter') @@ HashRecord
                                        /\ cursor_element' = HashCounter'
                             /\ challenge_path_acc' = challenge_path_acc + place_acc
                        ELSE /\ IF (<<cursor_element, (proof_elements[proof_index]), height>> \in DOMAIN HashRecord)
                                   THEN /\ cursor_element' = HashRecord[<<cursor_element, (proof_elements[proof_index]), height>>]
                                        /\ UNCHANGED << HashCounter, 
                                                        HashRecord >>
                                   ELSE /\ HashCounter' = HashCounter - 1
                                        /\ HashRecord' = (<<cursor_element, (proof_elements[proof_index]), height>> :> HashCounter') @@ HashRecord
                                        /\ cursor_element' = HashCounter'
                             /\ UNCHANGED challenge_path_acc
                  /\ place_acc' = place_acc * 2
                  /\ proof_index' = proof_index + 1
                  /\ height' = height + 1
                  /\ IF height' < Len(proof_elements)
                        THEN /\ pc' = [pc EXCEPT !["merkle tree"] = "ProofCheckLoop"]
                        ELSE /\ pc' = [pc EXCEPT !["merkle tree"] = "CheckRoot"]
                  /\ UNCHANGED << h1, h2, h3, h, input, row, rowSize, nextRow, 
                                  index, proof_element, root, challenge, 
                                  cursor_index, cursor_row, rows, proof_path, 
                                  proof_elements >>

CheckRoot == /\ pc["merkle tree"] = "CheckRoot"
             /\ Assert(cursor_element = root, 
                       "Failure of assertion at line 140, column 9.")
             /\ Assert(challenge_path_acc = challenge - 1, 
                       "Failure of assertion at line 141, column 9.")
             /\ pc' = [pc EXCEPT !["merkle tree"] = "IncrementChallenge"]
             /\ UNCHANGED << HashCounter, HashRecord, h1, h2, h3, h, input, 
                             row, rowSize, nextRow, index, proof_element, root, 
                             challenge, cursor_index, cursor_row, 
                             cursor_element, proof_index, challenge_path_acc, 
                             place_acc, rows, height, proof_path, 
                             proof_elements >>

IncrementChallenge == /\ pc["merkle tree"] = "IncrementChallenge"
                      /\ challenge' = challenge + 1
                      /\ IF challenge' <= Len(input)
                            THEN /\ pc' = [pc EXCEPT !["merkle tree"] = "MakeProof"]
                            ELSE /\ pc' = [pc EXCEPT !["merkle tree"] = "Done"]
                      /\ UNCHANGED << HashCounter, HashRecord, h1, h2, h3, h, 
                                      input, row, rowSize, nextRow, index, 
                                      proof_element, root, cursor_index, 
                                      cursor_row, cursor_element, proof_index, 
                                      challenge_path_acc, place_acc, rows, 
                                      height, proof_path, proof_elements >>

merkle_tree == BuildTree \/ RowLoop \/ HashRow \/ Advance \/ Repeat
                  \/ Proofs \/ MakeProof \/ S1 \/ ProofLoop \/ FinishProof
                  \/ CheckProof \/ ProofCheckLoop \/ CheckRoot
                  \/ IncrementChallenge

Next == test_hash \/ merkle_tree
           \/ (* Disjunct to prevent deadlock on termination *)
              ((\A self \in ProcSet: pc[self] = "Done") /\ UNCHANGED vars)

Spec == /\ Init /\ [][Next]_vars
        /\ WF_vars(merkle_tree)

Termination == <>(\A self \in ProcSet: pc[self] = "Done")

\* END TRANSLATION
=============================================================================
