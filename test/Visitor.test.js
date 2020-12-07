const app = require("../app");
const chai = require("chai");
const chaiHttp = require("chai-http");

const { expect } = chai;
chai.use(chaiHttp);

describe("Visitor:", () => {
  // Testing , that will change in POST Request
  let id = 0;

  describe("GET /api/Visitor", () => {
    it("GET All elements", (done) => {
      chai
        .request(app)
        .get("/api/Visitor/1")
        .end((err, res) => {
          expect(res).to.have.status(200);
          expect(res.body).to.have.a("array");
          expect(res.body.length > 0).to.equal(true);

          for (let i in res.body) {
            // Check keys
            expect(res.body[i]).to.have.property("id");
            expect(res.body[i]).to.have.property("Country");
            expect(res.body[i]).to.have.property("ip");
            expect(res.body[i]).to.have.property("Visit_Date");
            expect(res.body[i]).to.have.property("Count");
          }

          done();
        });
    });

    it("GET First element by :id", (done) => {
      chai
        .request(app)
        .get("/api/Visitor/1")
        .end((err, res) => {
          expect(res).to.have.status(200);
          expect(res.body).to.have.a("array");
          expect(res.body.length).to.equal(1);

          // Check keys
          expect(res.body[0]).to.have.property("id");
          expect(res.body[0]).to.have.property("Country");
          expect(res.body[0]).to.have.property("ip");
          expect(res.body[0]).to.have.property("Visit_Date");
          expect(res.body[0]).to.have.property("Count");

          // Check Data
          expect(res.body[0].id).to.equal(1);
          done();
        });
    });

    it("GET First element by params", (done) => {
      chai
        .request(app)
        .get("/api/Visitor?id=1")
        .end((err, res) => {
          expect(res).to.have.status(200);
          expect(res.body).to.have.a("array");
          expect(res.body.length).to.equal(1);

          // Check keys
          expect(res.body[0]).to.have.property("id");
          expect(res.body[0]).to.have.property("Country");
          expect(res.body[0]).to.have.property("ip");
          expect(res.body[0]).to.have.property("Visit_Date");
          expect(res.body[0]).to.have.property("Count");

          // Check Data
          expect(res.body[0].id).to.equal(1);
          done();
        });
    });

    it("GET not existing element", (done) => {
      chai
        .request(app)
        .get("/api/Visitor/-1")
        .end((err, res) => {
          expect(res).to.have.status(204);
          expect(res.body).to.be.empty;
          done();
        });
    });
  });

  describe("POST /api/Visitor", () => {
    const data = {
      Country: "UA",
      ip: "127.0.0.1",
      Visit_Date: "2000-12-07",
      Count: 3,
    };

    it("POST Request without Authorization", (done) => {
      chai
        .request(app)
        .post("/api/Visitor")
        .send(data)
        .end((err, res) => {
          expect(res).to.have.status(401);
          expect(res.body).to.have.property("message");
          expect(res.body.message).to.equal("Wrong Token declaration");
          done();
        });
    });

    it("POST Wrong Authorization Token", (done) => {
      chai
        .request(app)
        .post("/api/Visitor")
        .set("Authorization", "Bearer TEST")
        .send(data)
        .end((err, res) => {
          expect(res).to.have.status(403);
          expect(res.body).to.have.property("message");
          done();
        });
    });

    it("POST Successful Request", (done) => {
      chai
        .request(app)
        .post("/api/login")
        .send({ user: process.env.TEST_USER, pass: process.env.TEST_PASS })
        .end((err, res) => {
          expect(res).to.have.status(200);
          expect(res.body.accessToken).to.exist;

          chai
            .request(app)
            .post("/api/Visitor")
            .set("Authorization", "Bearer " + res.body.accessToken)
            .send(data)
            .end((err, res) => {
              expect(res).to.have.status(201);

              // Check keys
              expect(res.body).to.have.property("id");
              expect(res.body).to.have.property("Country");
              expect(res.body).to.have.property("ip");
              expect(res.body).to.have.property("Visit_Date");
              expect(res.body).to.have.property("Count");

              // Check Data
              id = res.body.id;
              expect(res.body.Country).to.equal(data.Country);
              expect(res.body.ip).to.equal(data.ip);
              expect(res.body.Visit_Date).to.equal(data.Visit_Date);
              expect(res.body.Count).to.equal(data.Count);
              done();
            });
        });
    });
  });

  describe("PUT /api/Visitor", () => {
    const data = {
      Visit_Date: "2020-12-07",
    };

    it("PUT Request without Authorization", (done) => {
      chai
        .request(app)
        .put("/api/Visitor")
        .send(data)
        .end((err, res) => {
          expect(res).to.have.status(401);
          expect(res.body).to.have.property("message");
          expect(res.body.message).to.equal("Wrong Token declaration");
          done();
        });
    });

    it("PUT Wrong Authorization Token", (done) => {
      chai
        .request(app)
        .put("/api/Visitor")
        .set("Authorization", "Bearer TEST")
        .send(data)
        .end((err, res) => {
          expect(res).to.have.status(403);
          expect(res.body).to.have.property("message");
          done();
        });
    });

    it("PUT Successful Request", (done) => {
      chai
        .request(app)
        .post("/api/login")
        .send({ user: process.env.TEST_USER, pass: process.env.TEST_PASS })
        .end((err, res) => {
          expect(res).to.have.status(200);
          expect(res.body.accessToken).to.exist;

          chai
            .request(app)
            .put(`/api/Visitor/${id}`)
            .set("Authorization", "Bearer " + res.body.accessToken)
            .send(data)
            .end((err, res) => {
              expect(res).to.have.status(204);
              expect(res.body).to.be.empty;
              done();
            });
        });
    });

    it("GET Update Element :id", (done) => {
      chai
        .request(app)
        .get(`/api/Visitor/${id}`)
        .end((err, res) => {
          expect(res).to.have.status(200);
          expect(res.body).to.have.a("array");
          expect(res.body.length).to.equal(1);

          // Check keys
          expect(res.body[0]).to.have.property("id");
          expect(res.body[0]).to.have.property("Country");
          expect(res.body[0]).to.have.property("ip");
          expect(res.body[0]).to.have.property("Visit_Date");
          expect(res.body[0]).to.have.property("Count");

          // Check Data
          expect(res.body[0].Curr_Date).to.equal(data.Curr_Date);
          done();
        });
    });

    it("PUT Not Modified Request", (done) => {
      chai
        .request(app)
        .post("/api/login")
        .send({ user: process.env.TEST_USER, pass: process.env.TEST_PASS })
        .end((err, res) => {
          expect(res).to.have.status(200);
          expect(res.body.accessToken).to.exist;

          chai
            .request(app)
            .put(`/api/Visitor/-1`)
            .set("Authorization", "Bearer " + res.body.accessToken)
            .send(data)
            .end((err, res) => {
              expect(res).to.have.status(304);
              expect(res.body).to.be.empty;
              done();
            });
        });
    });
  });

  describe("DELETE /api/Visitor", () => {
    it("DELETE Request without Authorization", (done) => {
      chai
        .request(app)
        .del("/api/Visitor")
        .end((err, res) => {
          expect(res).to.have.status(401);
          expect(res.body).to.have.property("message");
          expect(res.body.message).to.equal("Wrong Token declaration");
          done();
        });
    });

    it("DELETE Wrong Authorization Token", (done) => {
      chai
        .request(app)
        .del("/api/Visitor")
        .set("Authorization", "Bearer TEST")
        .end((err, res) => {
          expect(res).to.have.status(403);
          expect(res.body).to.have.property("message");
          done();
        });
    });

    it("DELETE Successful Request", (done) => {
      chai
        .request(app)
        .post("/api/login")
        .send({ user: process.env.TEST_USER, pass: process.env.TEST_PASS })
        .end((err, res) => {
          expect(res).to.have.status(200);
          expect(res.body.accessToken).to.exist;

          chai
            .request(app)
            .del(`/api/Visitor/${id}`)
            .set("Authorization", "Bearer " + res.body.accessToken)
            .end((err, res) => {
              expect(res).to.have.status(204);
              expect(res.body).to.be.empty;
              done();
            });
        });
    });

    it("GET Deleted Element :id", (done) => {
      chai
        .request(app)
        .get(`/api/Visitor/${id}`)
        .end((err, res) => {
          expect(res).to.have.status(204);
          expect(res.body).to.be.empty;
          done();
        });
    });

    it("DELETE Not Modified Request", (done) => {
      chai
        .request(app)
        .post("/api/login")
        .send({ user: process.env.TEST_USER, pass: process.env.TEST_PASS })
        .end((err, res) => {
          expect(res).to.have.status(200);
          expect(res.body.accessToken).to.exist;

          chai
            .request(app)
            .del(`/api/Visitor/-1`)
            .set("Authorization", "Bearer " + res.body.accessToken)
            .end((err, res) => {
              expect(res).to.have.status(304);
              expect(res.body).to.be.empty;
              done();
            });
        });
    });
  });
});
