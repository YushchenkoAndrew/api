const app = require("../app");
const chai = require("chai");
const chaiHttp = require("chai-http");

const { expect } = chai;
chai.use(chaiHttp);

describe("API Server", () => {
  it("Server is alive", (done) => {
    chai
      .request(app)
      .get("/api/ping")
      .end((err, res) => {
        expect(res).to.have.status(200);
        expect(res.body.message).to.equals("OK");
        done();
      });
  });

  it("Check Authorization", (done) => {
    chai
      .request(app)
      .post("/api/login")
      .send({ user: process.env.TEST_USER, pass: process.env.TEST_PASS })
      .end((err, res) => {
        expect(res).to.have.status(200);
        expect(res.body.accessToken).to.exist;
        chai
          .request(app)
          .get("/api/token")
          .set("Authorization", "Bearer " + res.body.accessToken)
          .end((err, res) => {
            expect(res).to.have.status(200);
            expect(res.body.message).to.equals("OK");
            done();
          });
      });
  });
});

// describe("GET /api/Views", () => {});
