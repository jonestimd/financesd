import React from 'react';
import {shallow} from 'enzyme';
import Table from './Table';
import HeaderRow from './HeaderRow';
import Row from './Row';
import {IColumn} from './Column';
import {mockSelectionHook} from 'src/test/mockHooks';
import TableHead from '@material-ui/core/TableHead';
import TableBody from '@material-ui/core/TableBody';

class TestRow {
    constructor(
        readonly id: string,
        readonly c1: string,
        readonly c2: string | number,
        readonly c3: string
    ) { }
}

describe('Table', () => {
    const columns: IColumn<TestRow>[] = [
        {key: 'first', className: 'col1', render: (row) => row.c1},
        {key: 'second', className: (row) => typeof row?.c2, render: (row) => row.c2, colspan: 2},
        {key: 'third', render: (row) => row.c3, header: (key) => key.toUpperCase()},
    ];
    const data = [
        new TestRow('r1', 'a1', 'a2', 'a3'),
        new TestRow('r2', 'b1', 20, 'b3'),
    ];

    it('tracks selected cell', () => {
        const selection = mockSelectionHook();

        const component = shallow(<Table columns={columns} data={data} />);

        const container = component.find('.scroll-container');
        expect(container).toExist();
        expect(container).toHaveProp('onKeyDown', selection.onKeyDown);
        expect(container).toHaveProp('onMouseDown', selection.onMouseDown);
        expect(container).toHaveProp('tabIndex', 0);
    });
    it('displays header row', () => {
        const component = shallow(<Table columns={columns} data={data} />);

        const header = component.find(TableHead);

        expect(header.find(HeaderRow)).toHaveProp('columns', columns);
    });
    it('populates table body using data', () => {
        const component = shallow(<Table columns={columns} data={data} />);

        const body = component.find(TableBody);

        const rows = body.find(Row);
        expect(rows).toHaveLength(data.length);
        expect(rows.at(0)).toHaveClassName('selected');
        expect(rows.at(0)).toHaveProp('row', data[0]);
        expect(rows.at(1)).not.toHaveClassName('selected');
        expect(rows.at(1)).toHaveProp('row', data[1]);
    });
    it('sets editing cell on row click', () => {
        const col = 1;
        const component = shallow(<Table columns={columns} data={data} />);

        component.find(Row).at(1).prop('onClick')?.(col);

        expect(component.find(Row).at(0)).toHaveProp('editCell', false);
        expect(component.find(Row).at(1)).toHaveProp('editCell', col);
    });
    it('ends editing on row commit', () => {
        const mockRef = {current: {focus: jest.fn()}};
        jest.spyOn(React, 'useRef').mockReturnValue(mockRef);
        const col = 1;
        const component = shallow(<Table columns={columns} data={data} />);
        component.find(Row).at(1).prop('onClick')?.(col);

        component.find(Row).at(1).prop('onCommit')?.();

        expect(component.find(Row).at(0)).toHaveProp('editCell', false);
        expect(component.find(Row).at(1)).toHaveProp('editCell', false);
        expect(mockRef.current.focus).toBeCalledTimes(1);
    });
    it('handles null ref (that should never happen)', () => {
        jest.spyOn(React, 'useRef').mockReturnValue({current: null});
        const col = 1;
        const component = shallow(<Table columns={columns} data={data} />);
        component.find(Row).at(1).prop('onClick')?.(col);

        component.find(Row).at(1).prop('onCommit')?.();

        expect(component.find(Row).at(0)).toHaveProp('editCell', false);
        expect(component.find(Row).at(1)).toHaveProp('editCell', false);
    });
});
