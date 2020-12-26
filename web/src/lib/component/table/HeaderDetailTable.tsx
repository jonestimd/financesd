import React from 'react';
import {observer} from 'mobx-react-lite';
import classNames from 'classnames';
import {HeaderRow, Row, IColumn, IRow} from './Table';
import IMixedRowTableModel from '../../model/IMixedRowTableModel';
import ScrollViewport, {IScrollableProps} from '../scroll/ScrollViewport';
import {useScroll} from '../scroll/scrollHooks';
import {useSelection} from '../scroll/selectionHooks';

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
    const scroll = useScroll<HTMLTableElement>(scrollOptions);
    // React.useEffect(() => {
    //     if (tableRef.current && model.groups.length > 0) gotoBottom();
    // }, [model, tableRef.current, rowHeight]);
    const startGroup = model.getGroupIndex(scroll.startRow);
    const leadingHeight = model.precedingRows[startGroup] * scroll.rowHeight;
    const height = model.rowCount * scroll.rowHeight + scroll.headerHeight;
    const selection = useSelection(scroll.startRow, model.rowCount, columns.length, Math.max(0, model.precedingRows[startGroup] - 1));
    return (
        <ScrollViewport onScroll={scroll.onScroll} onKeyDown={selection.onKeyDown} onMouseDown={selection.onMouseDown}>
            {({scrollHeight}: IScrollableProps) => {
                const endGroup = model.getGroupIndex(scroll.startRow + Math.ceil(scrollHeight / scroll.rowHeight) * 3);
                const trailingHeight = model.getRowsAfter(endGroup) * scroll.rowHeight;
                const rowClassNames = (index: number, subIndex = 0, classes?: string) => classNames(classes, {
                    even: (startGroup + index) % 2,
                    selected: model.precedingRows[startGroup + index] + subIndex === selection.row,
                });
                return (
                    <table ref={scroll.listRef} className={classNames('table header-detail', className)} style={{height}}>
                        <thead>
                            <HeaderRow columns={columns} />
                            <HeaderRow className='detail' columns={subColumns} />
                        </thead>
                        <tbody>
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
                            <tr className='prototype'><td>0</td></tr>
                        </tbody>
                    </table>
                );
            }}
        </ScrollViewport>
    );
});

export default HeaderDetailTable;
