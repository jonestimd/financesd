import React from 'react';
import {TableRow, TableCell} from '@material-ui/core';
import {translate} from '../../i18n/localize';
import {columnClasses, IColumn} from './Column';

export interface IHeaderProps<T> {
    className?: string;
    columns: IColumn<T>[];
}

const HeaderRow: React.FC<IHeaderProps<never>> = ({className, columns}) => {
    const classes = columnClasses(columns);
    return (
        <TableRow className={className}>
            {columns.map(({key, colspan, header = translate}, index) =>
                <TableCell key={key} className={classes[index]} colSpan={colspan}>{header(key)}</TableCell>
            )}
        </TableRow>
    );
};

export default HeaderRow;
