export default interface IMixedRowTableModel<T> {
    readonly groups: T[];

    /**
     * Get the index of the group containing a row.
     * @returns the group index and the number of rows before the group
     */
    getGroupIndex(rowIndex: number): [number, number];

    readonly rowCount: number;
}