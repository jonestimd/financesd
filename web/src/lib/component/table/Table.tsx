import React, {ReactNode} from 'react';
import classNames from 'classnames';
import {translate} from '../../i18n/localize';
import {useSelection} from './selection';
import ScrollViewport from '../ScrollViewport';

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

function evalSupplier<T>(supplier: ClassSupplier<T>, row?: T): string {
    return typeof supplier === 'function' ? supplier(row) : supplier;
}

function columnClasses<T>(columns: IColumn<T>[], selectedIndex: number = -1, row?: T): string[] {
    return columns.reduce(({index, classes}, {className, colspan = 1}) => {
        const selected = index <= selectedIndex && selectedIndex < index + colspan;
        index += colspan;
        return {index, classes: classes.concat(classNames(evalSupplier(className, row), {selected}))};
    }, {index: 0, classes: []}).classes;
}

interface IHeaderProps<T> {
    className?: string;
    columns: IColumn<T>[];
}

export const HeaderRow: React.FC<IHeaderProps<any>> = ({className, columns}) => {
    const classes = columnClasses(columns);
    return (
        <tr className={className}>
            {columns.map(({key, colspan, header = translate}, index) =>
                <th key={key} className={classes[index]} colSpan={colspan}>{header(key)}</th>
            )}
        </tr>
    );
};

interface IRowProps<T> extends IHeaderProps<T> {
    row: T;
    selection?: {
        row?: number;
        column?: number;
    };
}

export const Row: React.FC<IRowProps<any>> = <T extends IRow>({row, className, columns, selection = {}}: IRowProps<T>) => {
    const classes = columnClasses(columns, selection.column, row);
    return (
        <tr className={className}>
            {columns.map(({key, render, colspan}, index) =>
                <td key={key} className={classes[index]} colSpan={colspan}>{render(row)}</td>
            )}
        </tr>
    );
};

const rowClass = (index: number, selection: {row: number}) => classNames({
    even: index % 2,
    selected: index === selection.row
});

const Table: React.FC<ITableProps<any>> = <T extends IRow>({columns, data, className}: ITableProps<T>) => {
    const selection = useSelection(0, data.length, columns.length);
    return (
        <ScrollViewport onKeyDown={selection.onKeyDown} onMouseDown={selection.onMouseDown}>
            <table className={classNames('table', className)}>
                <thead><HeaderRow columns={columns}/></thead>
                <tbody>
                    {data.map((row, index) =>
                        <Row key={row.id} row={row} className={rowClass(index, selection)} columns={columns} selection={selection} />
                    )}
                </tbody>
            </table>
        </ScrollViewport>
    );
};

export default Table;