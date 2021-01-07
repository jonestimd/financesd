import {newSecurity} from "src/test/securityFactory";
import {SecurityModel} from "./SecurityModel";

describe('SecurityModel', () => {
    const security = newSecurity();

    describe('constructor', () => {
        it('populates security properties', () => {
            const model = new SecurityModel(security);

            expect(model).toEqual(security);
        });
    });
});
