import React from 'react';
import {shallow} from 'enzyme';
import GroupInput from './GroupInput';
import {newDetail} from 'src/test/detailFactory';
import Autocomplete, {AutocompleteRenderInputParams} from '@material-ui/lab/Autocomplete';
import {RootStore} from 'src/lib/store/RootStore';
import {newGroup} from 'src/test/groupFactory';
import {Icon} from '@material-ui/core';

describe('GroupInput', () => {
    const {groupStore} = new RootStore();
    const groups = [newGroup({name: 'Some grOUp'}), newGroup({name: 'anOther GrouP'})];
    groups.forEach((g) => groupStore['groupsById'].set(g.id, g));

    beforeEach(() => {
        jest.spyOn(React, 'useContext').mockReturnValue({groupStore});
    });
    it('displays list of accounts and categories', () => {
        const detail = newDetail();

        const component = shallow(<GroupInput detail={detail} />);

        expect(component.find(Autocomplete)).toHaveProp('options', groupStore.groups);
        expect(component.find(Autocomplete)).toHaveProp('value', null);
    });
    it('uses name for option labels', () => {
        const detail = newDetail();
        const component = shallow(<GroupInput detail={detail} />);

        const getOptionLabel = component.find(Autocomplete).prop('getOptionLabel')!;

        expect(getOptionLabel(groups[0])).toEqual(groups[0].name);
    });
    it('displays options containing text, ignoring case', () => {
        const detail = newDetail();
        const component = shallow(<GroupInput detail={detail} />);

        const filterOptions = component.find(Autocomplete).prop('filterOptions')!;

        expect(filterOptions(groups, {inputValue: 'group', getOptionLabel: () => ''})).toEqual(groups);
        expect(filterOptions(groups, {inputValue: 'some', getOptionLabel: () => ''})).toEqual(groups.slice(0,1));
    });
    it('displays selected group', () => {
        const detail = newDetail({transactionGroupId: groups[1].id});

        const component = shallow(<GroupInput detail={detail} />);

        expect(component.find(Autocomplete)).toHaveProp('value', groups[1]);
    });
    it('displays input with workspaces icon', () => {
        const detail = newDetail();
        const component = shallow(<GroupInput detail={detail} />);
        const renderInput = component.find(Autocomplete).prop('renderInput')!;

        const input = shallow(renderInput({} as AutocompleteRenderInputParams) as React.ReactElement);

        expect(input).toHaveProp('InputProps', expect.objectContaining({
            startAdornment: <Icon>workspaces</Icon>,
        }));
    });
});
