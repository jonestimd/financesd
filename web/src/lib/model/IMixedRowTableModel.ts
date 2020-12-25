export default interface IMixedRowTableModel<T> {
    readonly groups: T[];

    /**
     * Get the index of the group containing a row.
     */
    getGroupIndex(rowIndex: number): number;
    getRowsAfter(groupIndex: number): number;

    readonly precedingRows: number[];
    readonly rowCount: number;
}
