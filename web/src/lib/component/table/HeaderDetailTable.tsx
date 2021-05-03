import React from 'react';
import {observer} from 'mobx-react-lite';
import classNames from 'classnames';
import {IRow} from './Table';
import HeaderRow from './HeaderRow';
import Row from './Row';
import {IColumn} from './Column';
import IMixedRowTableModel from '../../model/IMixedRowTableModel';
import ScrollViewport, {IScrollableProps} from '../scroll/ScrollViewport';
import {useSelection} from '../scroll/tableSelectionHooks';
import {Table, TableRow, TableCell, TableHead, TableBody} from '@material-ui/core';

export interface IHeaderDetailTableProps<T, S> {
    className?: string;
    columns: IColumn<T>[];
    subColumns: IColumn<S>[];
    subrows: (row: T) => S[];
    model: IMixedRowTableModel<T>;
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
type TableType = React.FC<IHeaderDetailTableProps<any, any>>;

const scrollOptions = {
    prototypeSelector: 'tbody tr.prototype',
    defaultRowHeight: 22,
    headerSelector: 'thead',
    defaultHeaderHeight: 22,
};

const HeaderDetailTable: TableType = observer(<T extends IRow, S extends IRow>(props: IHeaderDetailTableProps<T, S>) => {
    const {columns, subColumns, model, subrows, className} = props;
    // React.useEffect(() => {
    //     if (tableRef.current && model.groups.length > 0) gotoBottom();
    // }, [model, tableRef.current, rowHeight]);
    const selection = useSelection<HTMLTableElement>({...scrollOptions, rows: model.rowCount, columns: columns.length});
    const startGroup = model.getGroupIndex(selection.scroll.startRow);
    const leadingHeight = model.precedingRows[startGroup] * selection.scroll.rowHeight;
    const height = model.rowCount * selection.scroll.rowHeight + selection.scroll.headerHeight;
    return (
        <ScrollViewport onScroll={selection.scroll.onScroll} onKeyDown={selection.onKeyDown} onMouseDown={selection.onMouseDown}>
            {({scrollHeight}: IScrollableProps) => {
                const endGroup = model.getGroupIndex(selection.scroll.endRow(scrollHeight));
                const trailingHeight = model.getRowsAfter(endGroup) * selection.scroll.rowHeight;
                const rowClassNames = (index: number, subIndex = 0, classes?: string) => {
                    return classNames(classes, {
                        odd: (startGroup + index) % 2,
                        selected: model.precedingRows[startGroup + index] + subIndex === selection.row,
                    });
                };
                return (
                    <Table ref={selection.scroll.listRef} className={classNames('table header-detail', className)} style={{height}}>
                        <TableHead>
                            <HeaderRow columns={columns} />
                            <HeaderRow className='detail' columns={subColumns} />
                        </TableHead>
                        <TableBody>
                            {leadingHeight > 0 ? <tr style={{height: leadingHeight}}><td /></tr> : null}
                            {model.groups.slice(startGroup, endGroup + 1).map((row, index) =>
                                <React.Fragment key={row.id}>
                                    <Row row={row} className={rowClassNames(index)} columns={columns} selection={selection} />
                                    {subrows(row).map((subrow, subIndex) =>
                                        <Row key={subrow.id} row={subrow} className={rowClassNames(index, subIndex + 1, 'detail')}
                                            columns={subColumns} selection={selection} />
                                    )}
                                </React.Fragment>
                            )}
                            {trailingHeight > 0 ? <tr style={{height: trailingHeight}}><td /></tr> : null}
                            <TableRow className='prototype'><TableCell>0</TableCell></TableRow>
                        </TableBody>
                    </Table>
                );
            }}
        </ScrollViewport>
    );
});

export default HeaderDetailTable;
