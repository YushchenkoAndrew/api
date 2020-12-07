const app = require("../app");
const chai = require("chai");
const chaiHttp = require("chai-http");

const { expect } = chai;
chai.use(chaiHttp);

describe("Basic API Server Request", () => {
  describe("GET /api", () => {
    it("Get Navigation url", (done) => {
      chai
        .request(app)
        .get("/api")
        .end((err, res) => {
          expect(res).to.have.status(200);
          expect(res.body).to.be.a("object");
          done();
        });
    });
  });

  describe("GET /api/tables", () => {
    it("Get list of tables", (done) => {
      chai
        .request(app)
        .get("/api/tables")
        .end((err, res) => {
          expect(res).to.have.status(200);
          expect(res.body).to.be.a("array");
          done();
        });
    });
  });

  describe("GET /api/ping", () => {
    it("Server is Alive", (done) => {
      chai
        .request(app)
        .get("/api/ping")
        .end((err, res) => {
          expect(res).to.have.status(200);
          expect(res.body.message).to.equals("OK");
          done();
        });
    });
  });

  describe("POST /api/login", () => {
    it("Authorization Complete", (done) => {
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
});
