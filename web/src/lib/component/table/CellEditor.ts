export interface ICellEditorProps<T> {
    value: T;
    onCommit: (value: T) => void;
}
