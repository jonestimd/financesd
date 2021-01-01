import {configure} from 'enzyme';
import Adapter from 'enzyme-adapter-react-16';
import 'jest-enzyme';
import {clearTimers} from 'mobx-react-lite';

configure({adapter: new Adapter()});

global.navigator = {
    language: 'en-US',
} as Navigator;

global.fetch = jest.fn();

afterAll(() => clearTimers());
