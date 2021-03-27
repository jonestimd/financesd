import React from 'react';
import {shallow} from 'enzyme';
import Table from './Table';
import HeaderRow from './HeaderRow';
import Row from './Row';
import {IColumn} from './Column';
import TableHead from '@material-ui/core/TableHead';
import TableBody from '@material-ui/core/TableBody';
import TextCellEditor from './TextCellEditor';
import SelectionModel from 'src/lib/model/SelectionModel';

class TestRow {
    constructor(
        readonly id: number,
        public c1: string,
        readonly c2: string | number,
        readonly c3: string
    ) { }
}

describe('Table', () => {
    const editor = {
        Component: TextCellEditor,
        getValue: (r: TestRow) => r.c1,
        setValue: (r: TestRow, value: string) => r.c1 = value,
    };
    const columns: IColumn<TestRow>[] = [
        {key: 'first', className: 'col1', render: (row) => row.c1, editor},
        {key: 'second', className: (row) => typeof row?.c2, render: (row) => row.c2, colspan: 2},
        {key: 'third', render: (row) => row.c3, header: (key) => key.toUpperCase()},
    ];
    const data = [
        new TestRow(1, 'a1', 'a2', 'a3'),
        new TestRow(2, 'b1', 20, 'b3'),
    ];

    it('tracks selected cell', () => {
        const selection = new SelectionModel({rows: 2, columns: columns.length});

        const component = shallow(<Table columns={columns} data={data} selection={selection} />);

        const container = component.find('.scroll-container');
        expect(container).toExist();
        expect(container).toHaveProp('onMouseDown', selection.onMouseDown);
    });
    it('displays header row', () => {
        const component = shallow(<Table columns={columns} data={data} />);

        const header = component.find(TableHead);

        expect(header.find(HeaderRow)).toHaveProp('columns', columns);
    });
    it('populates table body using data', () => {
        const selection = new SelectionModel({rows: 2, columns: columns.length});
        const component = shallow(<Table columns={columns} data={data} selection={selection} />);

        const body = component.find(TableBody);

        const rows = body.find(Row);
        expect(rows).toHaveLength(data.length);
        expect(rows.at(0)).toHaveClassName('selected');
        expect(rows.at(0)).toHaveProp('row', data[0]);
        expect(rows.at(1)).not.toHaveClassName('selected');
        expect(rows.at(1)).toHaveProp('row', data[1]);
    });
    it('uses selection.validate', () => {
        const selection = new SelectionModel({rows: 2, columns: columns.length});
        const validate = jest.fn();
        Object.assign(selection, {validate});
        const component = shallow(<Table columns={columns} data={data} selection={selection} />);

        component.find(Row).at(0).prop('validate')?.(data[0], 'key');

        expect(validate).toBeCalledWith(data[0], 'key');
    });
    it('uses selection.isEditable', () => {
        const selection = new SelectionModel({rows: 2, columns: columns.length});
        const isEditable = jest.fn();
        Object.assign(selection, {isEditable});
        const component = shallow(<Table columns={columns} data={data} selection={selection} />);

        component.find(Row).at(0).prop('isEditable')?.(data[0], 'key');

        expect(isEditable).toBeCalledWith(data[0], 'key');
    });
    it('sets editing cell on row click', () => {
        const selection = new SelectionModel({rows: 2, columns: columns.length});
        const component = shallow(<Table columns={columns} data={data} selection={selection} />);

        component.find(Row).at(1).prop('onClick')?.(1);

        expect(selection.editCell).toEqual({row: 1, column: 1});
    });
    it('ends editing on row commit', () => {
        const selection = new SelectionModel({rows: 2, columns: columns.length});
        jest.spyOn(selection, 'stopEditing');
        const component = shallow(<Table columns={columns} data={data} selection={selection} />);

        component.find(Row).at(1).prop('onCommit')?.();

        expect(selection.stopEditing).toBeCalledTimes(1);
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
    it('edits a cell', () => {
        const selection = new SelectionModel({rows: 2, columns: columns.length});
        selection.cell.row = 1;
        selection.editCell = selection.cell;

        const component = shallow(<Table columns={columns} data={data} selection={selection} />);

        const rows = component.find(Row);
        expect(rows.at(0)).toHaveProp('editCell', false);
        expect(rows.at(1)).toHaveProp('editCell', 0);
    });
    describe('onKeyDown', () => {
        it('does nothing if no selection', () => {
            const component = shallow(<Table columns={columns} data={data} />);

            const event = {ctrlKey: false, altKey: false, key: 'Left'};
            component.simulate('keydown', event);
        });
        it('calls selection hook for non-printable char on editable cell', () => {
            const selection = new SelectionModel({rows: 2, columns: columns.length});
            jest.spyOn(selection, 'onKeyDown').mockReturnValue();
            const component = shallow(<Table columns={columns} data={data} selection={selection} />);

            const event = {ctrlKey: false, altKey: false, key: 'Left'};
            component.simulate('keydown', event);

            expect(selection.onKeyDown).toBeCalledWith(event);
            expect(selection.editCell).toBeUndefined();
        });
        it('does not start editing for printable char on non-editable cell', () => {
            const selection = new SelectionModel({rows: 2, columns: columns.length});
            selection.cell.column = 1;
            jest.spyOn(selection, 'onKeyDown').mockReturnValue();
            const component = shallow(<Table columns={columns} data={data} selection={selection} />);

            component.simulate('keydown', {ctrlKey: false, altKey: false, key: 'x'});

            expect(selection.onKeyDown).not.toBeCalled();
            expect(selection.editCell).toBeUndefined();
        });
        const startEditTests = [
            {name: 'printable char', ctrlKey: false, altKey: false, key: 'x'},
            {name: 'F2', ctrlKey: false, altKey: false, key: 'F2'},
            {name: 'Backspace', ctrlKey: false, key: 'Backspace'},
            {name: 'alt+Backspace', altKey: false, key: 'Backspace'},
            {name: 'ctrl+Backspace', ctrlKey: true, key: 'Backspace'},
            {name: 'Delete', ctrlKey: false, key: 'Delete'},
            {name: 'alt+Delete', altKey: true, key: 'Delete'},
            {name: 'ctrl+Delete', ctrlKey: true, key: 'Delete'},
        ];
        startEditTests.forEach(({name, ctrlKey, altKey, key}) => {
            it(`starts editing for ${name} on editable cell`, () => {
                const selection = new SelectionModel({rows: 2, columns: columns.length});
                selection.cell.row = 1;
                jest.spyOn(selection, 'onKeyDown');
                const component = shallow(<Table columns={columns} data={data} selection={selection} />);

                component.simulate('keydown', {ctrlKey, altKey, key});

                expect(selection.onKeyDown).not.toBeCalled();
                expect(selection.editCell).toEqual({row: 1, column: 0});
            });
        });
    });
});
