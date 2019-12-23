import React, {ReactNode} from 'react';
import classnames from 'classnames';
import {translate} from '../i18n/localize';

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
        <table className={classnames('table', className)}>
            <thead>
                <tr>
                    {columns.map(({key, className: colClass, header = translate}) =>
                         <th key={key} className={getClassName(colClass)}>{header(key)}</th>
                    )}
                </tr>
            </thead>
            <tbody>
                {data.map(row => (
                    <tr key={row.id}>
                        {columns.map(column =>
                            <td key={column.key} className={getClassName(column.className, row)}>{column.render(row)}</td>
                        )}
                    </tr>
                ))}
            </tbody>
        </table>
    );
};

export default Table;