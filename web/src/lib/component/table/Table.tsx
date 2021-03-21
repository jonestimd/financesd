import React from 'react';
import classNames from 'classnames';
import MuiTable from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableHead from '@material-ui/core/TableHead';
import HeaderRow from './HeaderRow';
import Row, {IRowProps} from './Row';
import {IColumn} from './Column';
import {observer} from 'mobx-react-lite';
import {ICell} from '../../model/SelectionModel';

type ContainerRef<T> = ((instance: T | null) => void) | React.MutableRefObject<T | null> | null;

export interface ITableProps<T> {
    columns: IColumn<T>[];
    data: T[];
    className?: string;
    selection?: {
        rowClass?: (index: number) => string;
        cell: ICell;
        editCell?: ICell;
        setEditCell?: (cell?: ICell) => void;
        stopEditing?: () => void;
        isEditable?: (row: T, key: string) => boolean;
        onMouseDown: React.MouseEventHandler;
        onKeyDown: React.KeyboardEventHandler;
        validate?: (row: T, key: string) => string[];
        containerRef?: ContainerRef<HTMLDivElement>;
    };
}

export interface IRow {
    id: string;
}

const rowClass = (index: number) => ({odd: index % 2});

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
            else selection.onKeyDown(e);
        }
    };
    const isEditable = selection?.isEditable && selection.isEditable.bind(selection);
    const validate = selection?.validate && selection.validate.bind(selection);
    const onCommit = selection?.stopEditing && selection.stopEditing.bind(selection);
    const rowProps = (row: T, index: number): IRowProps<T> => ({
        row, columns, isEditable, validate, onCommit,
        className: classNames(selection?.rowClass?.(index), rowClass(index)),
        selection: selection?.cell,
        editCell: editCell?.row === index && editCell.column,
        onClick: (column) => selection?.setEditCell?.({row: index, column}),
    });
    return (
        <div className='scroll-container' onKeyDown={onKeyDown} onMouseDown={selection?.onMouseDown} ref={selection?.containerRef}>
            <MuiTable className={classNames('table', className)}>
                <TableHead><HeaderRow columns={columns} /></TableHead>
                <TableBody>{data.map((row, index) => <Row key={row.id} {...rowProps(row, index)} />)}</TableBody>
            </MuiTable>
        </div>
    );
};

export default observer(Table);
