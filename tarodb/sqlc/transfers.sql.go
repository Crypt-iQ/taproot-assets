// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: transfers.sql

package sqlc

import (
	"context"
	"database/sql"
	"time"
)

const applyPendingOutput = `-- name: ApplyPendingOutput :one
WITH spent_asset AS (
    SELECT genesis_id, version, asset_group_sig_id, script_version, lock_time,
           relative_lock_time
    FROM assets
    WHERE assets.asset_id = $6
)
INSERT INTO assets (
    genesis_id, version, asset_group_sig_id, script_version, lock_time,
    relative_lock_time, script_key_id, anchor_utxo_id, amount,
    split_commitment_root_hash, split_commitment_root_value
) VALUES (
    (SELECT genesis_id FROM spent_asset),
    (SELECT version FROM spent_asset),
    (SELECT asset_group_sig_id FROM spent_asset),
    (SELECT script_version FROM spent_asset),
    (SELECT lock_time FROM spent_asset),
    (SELECT relative_lock_time FROM spent_asset),
    $1, $2, $3, $4,
    $5
)
RETURNING asset_id
`

type ApplyPendingOutputParams struct {
	ScriptKeyID              int32
	AnchorUtxoID             sql.NullInt32
	Amount                   int64
	SplitCommitmentRootHash  []byte
	SplitCommitmentRootValue sql.NullInt64
	SpentAssetID             int32
}

func (q *Queries) ApplyPendingOutput(ctx context.Context, arg ApplyPendingOutputParams) (int32, error) {
	row := q.db.QueryRowContext(ctx, applyPendingOutput,
		arg.ScriptKeyID,
		arg.AnchorUtxoID,
		arg.Amount,
		arg.SplitCommitmentRootHash,
		arg.SplitCommitmentRootValue,
		arg.SpentAssetID,
	)
	var asset_id int32
	err := row.Scan(&asset_id)
	return asset_id, err
}

const deleteAssetWitnesses = `-- name: DeleteAssetWitnesses :exec
DELETE FROM asset_witnesses
WHERE asset_id = $1
`

func (q *Queries) DeleteAssetWitnesses(ctx context.Context, assetID int32) error {
	_, err := q.db.ExecContext(ctx, deleteAssetWitnesses, assetID)
	return err
}

const fetchTransferInputs = `-- name: FetchTransferInputs :many
SELECT input_id, anchor_point, asset_id, script_key, amount
FROM asset_transfer_inputs inputs
WHERE transfer_id = $1
ORDER BY input_id
`

type FetchTransferInputsRow struct {
	InputID     int32
	AnchorPoint []byte
	AssetID     []byte
	ScriptKey   []byte
	Amount      int64
}

func (q *Queries) FetchTransferInputs(ctx context.Context, transferID int32) ([]FetchTransferInputsRow, error) {
	rows, err := q.db.QueryContext(ctx, fetchTransferInputs, transferID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FetchTransferInputsRow
	for rows.Next() {
		var i FetchTransferInputsRow
		if err := rows.Scan(
			&i.InputID,
			&i.AnchorPoint,
			&i.AssetID,
			&i.ScriptKey,
			&i.Amount,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const fetchTransferOutputs = `-- name: FetchTransferOutputs :many
SELECT
    output_id, proof_suffix, amount, serialized_witnesses, script_key_local,
    split_commitment_root_hash, split_commitment_root_value, num_passive_assets,
    passive_assets_only,
    utxos.utxo_id AS anchor_utxo_id,
    utxos.outpoint AS anchor_outpoint,
    utxos.amt_sats AS anchor_value,
    utxos.merkle_root AS anchor_merkle_root,
    utxos.tapscript_sibling AS anchor_tapscript_sibling,
    utxo_internal_keys.raw_key AS internal_key_raw_key_bytes,
    utxo_internal_keys.key_family AS internal_key_family,
    utxo_internal_keys.key_index AS internal_key_index,
    script_keys.tweaked_script_key AS script_key_bytes,
    script_keys.tweak AS script_key_tweak,
    script_key AS script_key_id,
    script_internal_keys.raw_key AS script_key_raw_key_bytes,
    script_internal_keys.key_family AS script_key_family,
    script_internal_keys.key_index AS script_key_index
FROM asset_transfer_outputs outputs
JOIN managed_utxos utxos
  ON outputs.anchor_utxo = utxos.utxo_id
JOIN script_keys
  ON outputs.script_key = script_keys.script_key_id
JOIN internal_keys script_internal_keys
  ON script_keys.internal_key_id = script_internal_keys.key_id
JOIN internal_keys utxo_internal_keys
  ON utxos.internal_key_id = utxo_internal_keys.key_id
WHERE transfer_id = $1
ORDER BY output_id
`

type FetchTransferOutputsRow struct {
	OutputID                 int32
	ProofSuffix              []byte
	Amount                   int64
	SerializedWitnesses      []byte
	ScriptKeyLocal           bool
	SplitCommitmentRootHash  []byte
	SplitCommitmentRootValue sql.NullInt64
	NumPassiveAssets         int32
	PassiveAssetsOnly        bool
	AnchorUtxoID             int32
	AnchorOutpoint           []byte
	AnchorValue              int64
	AnchorMerkleRoot         []byte
	AnchorTapscriptSibling   []byte
	InternalKeyRawKeyBytes   []byte
	InternalKeyFamily        int32
	InternalKeyIndex         int32
	ScriptKeyBytes           []byte
	ScriptKeyTweak           []byte
	ScriptKeyID              int32
	ScriptKeyRawKeyBytes     []byte
	ScriptKeyFamily          int32
	ScriptKeyIndex           int32
}

func (q *Queries) FetchTransferOutputs(ctx context.Context, transferID int32) ([]FetchTransferOutputsRow, error) {
	rows, err := q.db.QueryContext(ctx, fetchTransferOutputs, transferID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FetchTransferOutputsRow
	for rows.Next() {
		var i FetchTransferOutputsRow
		if err := rows.Scan(
			&i.OutputID,
			&i.ProofSuffix,
			&i.Amount,
			&i.SerializedWitnesses,
			&i.ScriptKeyLocal,
			&i.SplitCommitmentRootHash,
			&i.SplitCommitmentRootValue,
			&i.NumPassiveAssets,
			&i.PassiveAssetsOnly,
			&i.AnchorUtxoID,
			&i.AnchorOutpoint,
			&i.AnchorValue,
			&i.AnchorMerkleRoot,
			&i.AnchorTapscriptSibling,
			&i.InternalKeyRawKeyBytes,
			&i.InternalKeyFamily,
			&i.InternalKeyIndex,
			&i.ScriptKeyBytes,
			&i.ScriptKeyTweak,
			&i.ScriptKeyID,
			&i.ScriptKeyRawKeyBytes,
			&i.ScriptKeyFamily,
			&i.ScriptKeyIndex,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertAssetTransfer = `-- name: InsertAssetTransfer :one
WITH target_txn(txn_id) AS (
    SELECT txn_id
    FROM chain_txns
    WHERE txid = $3
)
INSERT INTO asset_transfers (
    height_hint, anchor_txn_id, transfer_time_unix
) VALUES (
    $1, (SELECT txn_id FROM target_txn), $2
) RETURNING id
`

type InsertAssetTransferParams struct {
	HeightHint       int32
	TransferTimeUnix time.Time
	AnchorTxid       []byte
}

func (q *Queries) InsertAssetTransfer(ctx context.Context, arg InsertAssetTransferParams) (int32, error) {
	row := q.db.QueryRowContext(ctx, insertAssetTransfer, arg.HeightHint, arg.TransferTimeUnix, arg.AnchorTxid)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const insertAssetTransferInput = `-- name: InsertAssetTransferInput :exec
INSERT INTO asset_transfer_inputs (
    transfer_id, anchor_point, asset_id, script_key, amount
) VALUES (
    $1, $2, $3, $4, $5
)
`

type InsertAssetTransferInputParams struct {
	TransferID  int32
	AnchorPoint []byte
	AssetID     []byte
	ScriptKey   []byte
	Amount      int64
}

func (q *Queries) InsertAssetTransferInput(ctx context.Context, arg InsertAssetTransferInputParams) error {
	_, err := q.db.ExecContext(ctx, insertAssetTransferInput,
		arg.TransferID,
		arg.AnchorPoint,
		arg.AssetID,
		arg.ScriptKey,
		arg.Amount,
	)
	return err
}

const insertAssetTransferOutput = `-- name: InsertAssetTransferOutput :exec
INSERT INTO asset_transfer_outputs (
    transfer_id, anchor_utxo, script_key, script_key_local,
    amount, serialized_witnesses, split_commitment_root_hash,
    split_commitment_root_value, proof_suffix, num_passive_assets,
    passive_assets_only
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
`

type InsertAssetTransferOutputParams struct {
	TransferID               int32
	AnchorUtxo               int32
	ScriptKey                int32
	ScriptKeyLocal           bool
	Amount                   int64
	SerializedWitnesses      []byte
	SplitCommitmentRootHash  []byte
	SplitCommitmentRootValue sql.NullInt64
	ProofSuffix              []byte
	NumPassiveAssets         int32
	PassiveAssetsOnly        bool
}

func (q *Queries) InsertAssetTransferOutput(ctx context.Context, arg InsertAssetTransferOutputParams) error {
	_, err := q.db.ExecContext(ctx, insertAssetTransferOutput,
		arg.TransferID,
		arg.AnchorUtxo,
		arg.ScriptKey,
		arg.ScriptKeyLocal,
		arg.Amount,
		arg.SerializedWitnesses,
		arg.SplitCommitmentRootHash,
		arg.SplitCommitmentRootValue,
		arg.ProofSuffix,
		arg.NumPassiveAssets,
		arg.PassiveAssetsOnly,
	)
	return err
}

const insertPassiveAsset = `-- name: InsertPassiveAsset :exec
WITH target_asset(asset_id) AS (
    SELECT assets.asset_id
    FROM assets
        JOIN genesis_assets
            ON assets.genesis_id = genesis_assets.gen_asset_id
        JOIN managed_utxos utxos
            ON assets.anchor_utxo_id = utxos.utxo_id
        JOIN script_keys
            ON assets.script_key_id = script_keys.script_key_id
    WHERE genesis_assets.asset_id = $6
        AND utxos.outpoint = $7
        AND script_keys.tweaked_script_key = $3
)
INSERT INTO passive_assets (
    asset_id, transfer_id, new_anchor_utxo, script_key, new_witness_stack,
    new_proof
) VALUES (
    (SELECT asset_id FROM target_asset), $1, $2,
    $3, $4, $5
)
`

type InsertPassiveAssetParams struct {
	TransferID      int32
	NewAnchorUtxo   int32
	ScriptKey       []byte
	NewWitnessStack []byte
	NewProof        []byte
	AssetGenesisID  []byte
	PrevOutpoint    []byte
}

func (q *Queries) InsertPassiveAsset(ctx context.Context, arg InsertPassiveAssetParams) error {
	_, err := q.db.ExecContext(ctx, insertPassiveAsset,
		arg.TransferID,
		arg.NewAnchorUtxo,
		arg.ScriptKey,
		arg.NewWitnessStack,
		arg.NewProof,
		arg.AssetGenesisID,
		arg.PrevOutpoint,
	)
	return err
}

const insertReceiverProofTransferAttempt = `-- name: InsertReceiverProofTransferAttempt :exec
INSERT INTO receiver_proof_transfer_attempts (
    proof_locator_hash, time_unix
) VALUES (
    $1, $2
)
`

type InsertReceiverProofTransferAttemptParams struct {
	ProofLocatorHash []byte
	TimeUnix         time.Time
}

func (q *Queries) InsertReceiverProofTransferAttempt(ctx context.Context, arg InsertReceiverProofTransferAttemptParams) error {
	_, err := q.db.ExecContext(ctx, insertReceiverProofTransferAttempt, arg.ProofLocatorHash, arg.TimeUnix)
	return err
}

const queryAssetTransfers = `-- name: QueryAssetTransfers :many
SELECT
    id, height_hint, txns.txid, transfer_time_unix
FROM asset_transfers transfers
JOIN chain_txns txns
    ON transfers.anchor_txn_id = txns.txn_id
WHERE ($1 = false OR $1 IS NULL OR
    (CASE WHEN txns.block_hash IS NULL THEN true ELSE false END) = $1)

AND (txns.txid = $2 OR
    $2 IS NULL)
ORDER BY transfer_time_unix
`

type QueryAssetTransfersParams struct {
	UnconfOnly   interface{}
	AnchorTxHash []byte
}

type QueryAssetTransfersRow struct {
	ID               int32
	HeightHint       int32
	Txid             []byte
	TransferTimeUnix time.Time
}

// We'll use this clause to filter out for only transfers that are
// unconfirmed. But only if the unconf_only field is set.
// Here we have another optional query clause to select a given transfer
// based on the anchor_tx_hash, but only if it's specified.
func (q *Queries) QueryAssetTransfers(ctx context.Context, arg QueryAssetTransfersParams) ([]QueryAssetTransfersRow, error) {
	rows, err := q.db.QueryContext(ctx, queryAssetTransfers, arg.UnconfOnly, arg.AnchorTxHash)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []QueryAssetTransfersRow
	for rows.Next() {
		var i QueryAssetTransfersRow
		if err := rows.Scan(
			&i.ID,
			&i.HeightHint,
			&i.Txid,
			&i.TransferTimeUnix,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const queryPassiveAssets = `-- name: QueryPassiveAssets :many
SELECT passive.asset_id, passive.new_anchor_utxo, passive.script_key,
       passive.new_witness_stack, passive.new_proof,
       genesis_assets.asset_id AS genesis_id
FROM passive_assets as passive
    JOIN assets
        ON passive.asset_id = assets.asset_id
    JOIN genesis_assets
        ON assets.genesis_id = genesis_assets.gen_asset_id
WHERE passive.transfer_id = $1
`

type QueryPassiveAssetsRow struct {
	AssetID         int32
	NewAnchorUtxo   int32
	ScriptKey       []byte
	NewWitnessStack []byte
	NewProof        []byte
	GenesisID       []byte
}

func (q *Queries) QueryPassiveAssets(ctx context.Context, transferID int32) ([]QueryPassiveAssetsRow, error) {
	rows, err := q.db.QueryContext(ctx, queryPassiveAssets, transferID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []QueryPassiveAssetsRow
	for rows.Next() {
		var i QueryPassiveAssetsRow
		if err := rows.Scan(
			&i.AssetID,
			&i.NewAnchorUtxo,
			&i.ScriptKey,
			&i.NewWitnessStack,
			&i.NewProof,
			&i.GenesisID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const queryReceiverProofTransferAttempt = `-- name: QueryReceiverProofTransferAttempt :many
SELECT time_unix
FROM receiver_proof_transfer_attempts
WHERE proof_locator_hash = $1
ORDER BY time_unix DESC
`

func (q *Queries) QueryReceiverProofTransferAttempt(ctx context.Context, proofLocatorHash []byte) ([]time.Time, error) {
	rows, err := q.db.QueryContext(ctx, queryReceiverProofTransferAttempt, proofLocatorHash)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []time.Time
	for rows.Next() {
		var time_unix time.Time
		if err := rows.Scan(&time_unix); err != nil {
			return nil, err
		}
		items = append(items, time_unix)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const reAnchorPassiveAssets = `-- name: ReAnchorPassiveAssets :exec
UPDATE assets
SET anchor_utxo_id = $1
WHERE asset_id = $2
`

type ReAnchorPassiveAssetsParams struct {
	NewAnchorUtxoID sql.NullInt32
	AssetID         int32
}

func (q *Queries) ReAnchorPassiveAssets(ctx context.Context, arg ReAnchorPassiveAssetsParams) error {
	_, err := q.db.ExecContext(ctx, reAnchorPassiveAssets, arg.NewAnchorUtxoID, arg.AssetID)
	return err
}
