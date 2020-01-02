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

export function getClassName<T>(className: ClassSupplier<T>, row?: T, selected?: boolean) {
    return classNames(typeof className === 'function' ? className(row) : className, {selected});
}

interface IHeaderProps<T> {
    className?: string;
    columns: IColumn<T>[];
}

export const HeaderRow: React.FC<IHeaderProps<any>> = ({className, columns}) => {
    return (
        <tr className={className}>
            {columns.map(({key, className: colClass, colspan, header = translate}) =>
                <th key={key} className={getClassName(colClass)} colSpan={colspan}>{header(key)}</th>
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
    return (
        <tr className={className}>
            {columns.map(({key, className: colClass, render, colspan}, index) =>
                <td key={key} className={getClassName(colClass, row, index === selection.column)}
                    colSpan={colspan}>{render(row)}</td>
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