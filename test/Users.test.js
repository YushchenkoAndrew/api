const app = require("../app");
const chai = require("chai");
const chaiHttp = require("chai-http");

const { expect } = chai;
chai.use(chaiHttp);

describe("Users:", () => {
  describe("GET /api/Users", () => {
    it("GET Forbidden Data", (done) => {
      chai
        .request(app)
        .get("/api/Users")
        .end((err, res) => {
          expect(res).to.have.status(404);
          expect(res.body).to.have.property("message");
          done();
        });
    });
  });
});
