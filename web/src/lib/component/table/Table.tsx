import React, {useRef, useState} from 'react';
import classNames from 'classnames';
import {useSelection} from '../scroll/selectionHooks';
import MuiTable from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableHead from '@material-ui/core/TableHead';
import HeaderRow from './HeaderRow';
import Row from './Row';
import {IColumn} from './Column';

export interface ITableProps<T> {
    columns: IColumn<T>[];
    data: T[];
    className?: string;
}

export interface IRow {
    id: string;
}

const rowClass = (index: number, selection: {row: number}) => classNames({
    odd: index % 2,
    selected: index === selection.row,
});

export const selectionOptions = {
    rowSelector: 'tbody tr.MuiTableRow-root',
    headerSelector: 'thead',
};

const startEdit = (e: React.KeyboardEvent<HTMLElement>) => {
    return !e.ctrlKey && !e.altKey && (e.key.length === 1 || e.key === 'F2') || e.key === 'Delete' || e.key === 'Backspace';
};

const Table = <T extends IRow>({columns, data, className}: ITableProps<T>) => {
    const selection = useSelection({rows: data.length, columns: columns.length, ...selectionOptions});
    const [editCell, setEditCell] = useState<{row: number, column: number} | undefined>(undefined);
    const containerRef = useRef<HTMLDivElement>(null);
    const stopEditing = () => {
        setEditCell(undefined);
        containerRef.current?.focus();
    };
    const onKeyDown = (e: React.KeyboardEvent<HTMLElement>) => {
        if (startEdit(e)) {
            if (columns[selection.column].editor) {
                const {row, column} = selection;
                setEditCell({row, column});
            }
        }
        else selection.onKeyDown(e);
    };
    return (
        <div className='scroll-container' onKeyDown={onKeyDown} onMouseDown={selection.onMouseDown} tabIndex={0} ref={containerRef}>
            <MuiTable className={classNames('table', className)}>
                <TableHead><HeaderRow columns={columns} /></TableHead>
                <TableBody>
                    {data.map((row, index) =>
                        <Row<T> key={row.id} row={row} className={rowClass(index, selection)} columns={columns} selection={selection}
                            onClick={(column) => setEditCell({row: index, column})}
                            editCell={editCell?.row === index && editCell.column}
                            onCommit={stopEditing} />
                    )}
                </TableBody>
            </MuiTable>
        </div>
    );
};

export default Table;
