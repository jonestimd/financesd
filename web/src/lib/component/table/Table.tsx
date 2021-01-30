import React, {ReactNode} from 'react';
import classNames from 'classnames';
import {translate} from '../../i18n/localize';
import {useSelection} from '../scroll/selectionHooks';
import MuiTable from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import TableCell from '@material-ui/core/TableCell';

type ClassSupplier<T> = string | ((row?: T) => string);

export interface IColumn<T> {
    key: string;
    className?: ClassSupplier<T>;
    colspan?: number;
    header?: (key: string) => ReactNode;
    render: (row: T) => React.ReactNode;
}

export interface ITableProps<T> {
    columns: IColumn<T>[];
    data: T[];
    className?: string;
}

export interface IRow {
    id: string;
}

function evalSupplier<T>(supplier?: ClassSupplier<T>, row?: T): string {
    return typeof supplier === 'function' ? supplier(row) : supplier ?? '';
}

function columnClasses<T>(columns: IColumn<T>[], selectedIndex = -1, row?: T): string[] {
    return columns.reduce(({index, classes}, {className, colspan = 1}) => {
        const selected = index <= selectedIndex && selectedIndex < index + colspan;
        index += colspan;
        return {index, classes: classes.concat(classNames(evalSupplier(className, row), {selected}))};
    }, {index: 0, classes: [] as string[]}).classes;
}

interface IHeaderProps<T> {
    className?: string;
    columns: IColumn<T>[];
}

export const HeaderRow = <T extends unknown>({className, columns}: IHeaderProps<T>) => {
    const classes = columnClasses(columns);
    return (
        <TableRow className={className}>
            {columns.map(({key, colspan, header = translate}, index) =>
                <TableCell key={key} className={classes[index]} colSpan={colspan}>{header(key)}</TableCell>
            )}
        </TableRow>
    );
};

interface IRowProps<T> extends IHeaderProps<T> {
    row: T;
    selection?: {
        row?: number;
        column?: number;
    };
}

export const Row = <T extends IRow>({row, className, columns, selection = {}}: IRowProps<T>) => {
    const classes = columnClasses(columns, selection.column, row);
    return (
        <TableRow className={className}>
            {columns.map(({key, render, colspan}, index) =>
                <TableCell key={key} className={classes[index]} colSpan={colspan}>{render(row)}</TableCell>
            )}
        </TableRow>
    );
};

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
