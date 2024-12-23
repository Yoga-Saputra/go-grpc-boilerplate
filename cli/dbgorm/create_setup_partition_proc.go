package dbgorm

// Will create or replace DB Procedure of Setup Partition process.
func CreateSetupPartitionProc(a *AdvanceDBMigrate) error {
	if err := a.Tx.
		Exec(`
			CREATE OR REPLACE PROCEDURE setup_partition_proc(
				table_parent TEXT,
				partition_start_datetime VARCHAR(12)	
			)
			LANGUAGE plpgsql AS
			$$ BEGIN
				-- 	Create parent partition
				PERFORM partman.create_parent(
					p_parent_table => table_parent,
					p_control => 'date',
					p_type => 'native',
					p_interval => 'daily',
					p_start_partition => partition_start_datetime,
					p_premake => 90
				);
			
				-- 	Update partition config
				UPDATE partman.part_config
				SET infinite_time_partitions = true,
					retention                = '3 months',
					retention_keep_table     = false,
					retention_keep_index = false
				WHERE parent_table = table_parent;
			END $$;
		`).
		Error; err != nil {
		return err
	}

	return nil
}
