import React from 'react';
import classNames from 'classnames';
import {useSelection} from '../scroll/selectionHooks';
import MuiTable from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableHead from '@material-ui/core/TableHead';
import {IColumn, HeaderRow, Row} from './Row';

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

const Table = <T extends IRow>({columns, data, className}: ITableProps<T>) => {
    const selection = useSelection({rows: data.length, columns: columns.length, ...selectionOptions});
    return (
        <div className='scroll-container' onKeyDown={selection.onKeyDown} onMouseDown={selection.onMouseDown} tabIndex={0}>
            <MuiTable className={classNames('table', className)}>
                <TableHead><HeaderRow columns={columns} /></TableHead>
                <TableBody>
                    {data.map((row, index) =>
                        <Row key={row.id} row={row} className={rowClass(index, selection)} columns={columns} selection={selection} />
                    )}
                </TableBody>
            </MuiTable>
        </div>
    );
};

export default Table;
