import React, {ReactNode} from 'react';
import classNames from 'classnames';
import {translate} from '../i18n/localize';
import ScrollViewport from './ScrollViewport';

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

export function getClassName<T>(className: ClassSupplier<T>, row?: T) {
    return typeof className === 'function' ? className(row) : className;
}

const Table: React.FC<ITableProps<any>> = <T extends IRow>({columns, data, className}: ITableProps<T>) => {
    return (
        <ScrollViewport>
            <table className={classNames('table', className)}>
                <thead>
                    <tr>
                        {columns.map(({key, className: colClass, header = translate}) =>
                            <th key={key} className={getClassName(colClass)}>{header(key)}</th>
                        )}
                    </tr>
                </thead>
                <tbody>
                    {data.map((row, index) => (
                        <tr key={row.id} className={classNames({even: index % 2})}>
                            {columns.map(({key, className: colClass, render}) =>
                                <td key={key} className={getClassName(colClass, row)}>{render(row)}</td>
                            )}
                        </tr>
                    ))}
                </tbody>
            </table>
        </ScrollViewport>
    );
};

export default Table;