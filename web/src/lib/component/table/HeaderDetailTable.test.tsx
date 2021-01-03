import React from 'react';
import {mount, shallow} from 'enzyme';
import HeaderDetailTable from './HeaderDetailTable';
import IMixedRowTableModel from 'src/lib/model/IMixedRowTableModel';
import {sortedIndex} from 'lodash';
import {HeaderRow, IColumn, Row} from './Table';
import {mockHooks} from 'src/test/mockHooks';
import ScrollViewport from '../scroll/ScrollViewport';
import TableHead from '@material-ui/core/TableHead';
import TableBody from '@material-ui/core/TableBody';

let nextId = 1;

class TestData {
    readonly id: string;

    constructor(readonly c1: string, readonly c2: string | number, readonly c3: string) {
        this.id = `${nextId++}`;
    }
}

class TestRow extends TestData {
    constructor(c1: string, c2: string | number, c3: string, readonly details: TestData[] = []) {
        super(c1, c2, c3);
    }
}

class TestModel implements IMixedRowTableModel<TestRow> {
    readonly precedingRows: number[];
    readonly rowCount: number;

    constructor(readonly groups: TestRow[]) {
        this.precedingRows = this.groups.reduce((rows, group) => {
            const last = rows[rows.length - 1];
            return rows.concat(last + 1 + group.details.length);
        }, [0]);
        this.rowCount = this.precedingRows[this.groups.length];
    }

    getGroupIndex(rowIndex: number): number {
        const index = sortedIndex(this.precedingRows, rowIndex);
        return this.precedingRows[index] === rowIndex ? index : index - 1;
    }

    getRowsAfter(groupIndex: number): number {
        return this.rowCount - this.precedingRows[groupIndex + 1];
    }
}

const columns: IColumn<TestData>[] = [
    {key: 'first', className: 'col1', render: (row) => row.c1},
    {key: 'second', className: (row) => typeof row?.c2, render: (row) => row.c2, colspan: 2},
    {key: 'third', render: (row) => row.c3, header: (key) => key.toUpperCase()},
];

const subColumns: IColumn<TestData>[] = [
    {key: 'Beginning', className: 'col1', render: (row) => row.c1},
    {key: 'Middle', className: (row) => typeof row?.c2, render: (row) => row.c2, colspan: 2},
    {key: 'End', render: (row) => row.c3, header: (key) => key.toUpperCase()},
];

describe('HeaderDetailTable', () => {
    const data = [
        new TestRow('a1', 'a2', 'a3', [new TestData('x1', 'x2', 'x3'), new TestData('y1', 200, 'y3')]),
        new TestRow('b1', 20, 'b3', [new TestData('z1', 'z2', 'z3')]),
    ];
    const props = {
        data,
        columns,
        subColumns,
        model: new TestModel(data),
        subrows: (row: TestRow) => row.details,
    };
    const scrollHeight = 100;
    const rowHeight = 24;

    it('uses ScrollViewport', () => {
        const {scroll, selection} = mockHooks();

        const component = shallow(<HeaderDetailTable {...props} />);

        const viewport = component.find(ScrollViewport);
        expect(viewport).toHaveProp('onScroll', scroll.onScroll);
        expect(viewport).toHaveProp('onKeyDown', selection.onKeyDown);
        expect(viewport).toHaveProp('onMouseDown', selection.onMouseDown);
    });
    it('displays a table with headers', () => {
        const {scroll} = mockHooks();
        const component = shallow(<HeaderDetailTable {...props} />);

        const table = mount(component.find(ScrollViewport).prop('children')({scrollHeight}) as React.ReactElement);

        expect(scroll.listRef.current).toBeInstanceOf(HTMLTableElement);
        const headerRows = table.find(TableHead).find(HeaderRow);
        expect(headerRows).toHaveLength(2);
        expect(headerRows.at(0)).toHaveProp('columns', props.columns);
        expect(headerRows.at(1)).toHaveProp('columns', props.subColumns);
        expect(headerRows.at(1)).toHaveClassName('detail');
    });
    it('displays header and detail rows', () => {
        mockHooks();
        const component = shallow(<HeaderDetailTable {...props} />);

        const table = shallow(component.find(ScrollViewport).prop('children')({scrollHeight}) as React.ReactElement);

        const bodyRows = table.find(TableBody).find(Row);
        expect(bodyRows).toHaveLength(5);
        bodyRows.slice(0, 3).forEach((row) => expect(row).not.toHaveClassName('odd'));
        bodyRows.slice(3).forEach((row) => expect(row).toHaveClassName('odd'));
        expect(bodyRows.at(0)).toHaveClassName('selected');
        bodyRows.slice(1).forEach((row) => expect(row).not.toHaveClassName('selected'));
        expect(bodyRows.at(0)).not.toHaveClassName('detail');
        bodyRows.slice(1, 3).forEach((row) => expect(row).toHaveClassName('detail'));
        expect(bodyRows.at(3)).not.toHaveClassName('detail');
        expect(bodyRows.at(4)).toHaveClassName('detail');
        expect(bodyRows.at(0)).toHaveProp('row', props.data[0]);
        expect(bodyRows.at(1)).toHaveProp('row', props.data[0].details[0]);
        expect(bodyRows.at(2)).toHaveProp('row', props.data[0].details[1]);
        expect(bodyRows.at(3)).toHaveProp('row', props.data[1]);
        expect(bodyRows.at(4)).toHaveProp('row', props.data[1].details[0]);
    });
    it('displays leading filler', () => {
        const model = new TestModel(new Array(100).fill(null).map(() => new TestRow('b1', 20, 'b3', [new TestData('z1', 'z2', 'z3')])));
        const startRow = 75;
        mockHooks({startRow, rowHeight: 24});
        const component = shallow(<HeaderDetailTable {...props} model={model} />);

        const table = shallow(component.find(ScrollViewport).prop('children')({scrollHeight}) as React.ReactElement);

        expect(table.find('tr').first()).toHaveStyle({height: rowHeight * model.precedingRows[model.getGroupIndex(startRow)]});
    });
    it('displays trailing filler', () => {
        const model = new TestModel(new Array(100).fill(null).map(() => new TestRow('b1', 20, 'b3', [new TestData('z1', 'z2', 'z3')])));
        const endRow = 15;
        const {scroll} = mockHooks({rowHeight});
        scroll.endRow.mockReturnValue(endRow);
        const component = shallow(<HeaderDetailTable {...props} model={model} />);

        const table = shallow(component.find(ScrollViewport).prop('children')({scrollHeight}) as React.ReactElement);

        expect(table.find('tr').last()).toHaveStyle({height: rowHeight * model.getRowsAfter(model.getGroupIndex(endRow))});
    });
});
