import React from 'react';
import classNames from 'classnames';
import MuiTable from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableHead from '@material-ui/core/TableHead';
import HeaderRow from './HeaderRow';
import Row from './Row';
import {IColumn} from './Column';
import {observer} from 'mobx-react-lite';
import {ICell} from '../../model/SelectionModel';

type ContainerRef<T> = ((instance: T | null) => void) | React.MutableRefObject<T | null> | null;

export interface ITableProps<T> {
    columns: IColumn<T>[];
    data: T[];
    className?: string;
    selection?: {
        cell: ICell;
        editCell?: ICell;
        setEditCell?: (cell?: ICell) => void;
        stopEditing?: () => void;
        onMouseDown: React.MouseEventHandler;
        onKeyDown: React.KeyboardEventHandler;
        containerRef?: ContainerRef<HTMLDivElement>;
    };
}

export interface IRow {
    id: string;
}

const rowClass = (index: number, selection?: {row: number}) => classNames({
    odd: index % 2,
    selected: index === selection?.row,
});

const startEdit = (e: React.KeyboardEvent<HTMLElement>) => {
    return !e.ctrlKey && !e.altKey && (e.key.length === 1 || e.key === 'F2') || e.key === 'Delete' || e.key === 'Backspace';
};

const Table = <T extends IRow>({columns, data, className, selection}: ITableProps<T>) => {
    const editCell = selection?.editCell;
    const onKeyDown = (e: React.KeyboardEvent<HTMLElement>) => {
        if (selection) {
            if (selection.setEditCell && startEdit(e)) {
                if (columns[selection.cell.column].editor) {
                    selection.setEditCell(selection.cell);
                }
            }
            else selection.onKeyDown?.(e);
        }
    };
    return (
        <div className='scroll-container' onKeyDown={onKeyDown} onMouseDown={selection?.onMouseDown} ref={selection?.containerRef}>
            <MuiTable className={classNames('table', className)}>
                <TableHead><HeaderRow columns={columns} /></TableHead>
                <TableBody>
                    {data.map((row, index) =>
                        <Row<T> key={row.id} row={row} className={rowClass(index, selection?.cell)} columns={columns} selection={selection?.cell}
                            onClick={(column) => selection?.setEditCell?.({row: index, column})}
                            editCell={editCell?.row === index && editCell.column}
                            onCommit={selection?.stopEditing} />
                    )}
                </TableBody>
            </MuiTable>
        </div>
    );
};

export default observer(Table);
