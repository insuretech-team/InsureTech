////
////  FaceListViewController.swift
////  InsureTech
////
////  Created by LifeplusBD on 5/1/26.
////
//
//import UIKit
//
//
//final class FaceListViewController: UIViewController {
//
//    private var faceList: [FaceEntry] = []
//
//    private let tableView = UITableView()
//    private let captureButton = UIButton(type: .system)
//    private let verifyButton = UIButton(type: .system)
//
//    override func viewDidLoad() {
//        super.viewDidLoad()
//        view.backgroundColor = .white
//        setupUI()
//    }
//
//    private func setupUI() {
//        tableView.dataSource = self
//        tableView.register(UITableViewCell.self, forCellReuseIdentifier: "cell")
//        tableView.translatesAutoresizingMaskIntoConstraints = false
//
//        captureButton.setTitle("Capture", for: .normal)
//        verifyButton.setTitle("Verify", for: .normal)
//
//        captureButton.addTarget(self, action: #selector(captureTapped), for: .touchUpInside)
//        verifyButton.addTarget(self, action: #selector(verifyTapped), for: .touchUpInside)
//
//        let buttonStack = UIStackView(arrangedSubviews: [captureButton, verifyButton])
//        buttonStack.axis = .horizontal
//        buttonStack.spacing = 20
//        buttonStack.translatesAutoresizingMaskIntoConstraints = false
//
//        view.addSubview(tableView)
//        view.addSubview(buttonStack)
//
//        NSLayoutConstraint.activate([
//            tableView.topAnchor.constraint(equalTo: view.safeAreaLayoutGuide.topAnchor),
//            tableView.leadingAnchor.constraint(equalTo: view.leadingAnchor),
//            tableView.trailingAnchor.constraint(equalTo: view.trailingAnchor),
//
//            buttonStack.topAnchor.constraint(equalTo: tableView.bottomAnchor, constant: 10),
//            buttonStack.centerXAnchor.constraint(equalTo: view.centerXAnchor),
//            buttonStack.bottomAnchor.constraint(equalTo: view.safeAreaLayoutGuide.bottomAnchor, constant: -10)
//        ])
//    }
//
//    @objc private func captureTapped() {
//        let captureVC = FaceCaptureViewController()
//        captureVC.onFaceCaptured = { [weak self] name, image in
//            self?.faceList.append(FaceEntry(name: name, image: image))
//            self?.tableView.reloadData()
//        }
//        present(captureVC, animated: true)
//    }
//
//    @objc private func verifyTapped() {
//        // Push verify VC (optional)
//    }
//}
//
//extension FaceListViewController: UITableViewDataSource {
//    func tableView(_ tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
//        faceList.count
//    }
//
//    func tableView(_ tableView: UITableView, cellForRowAt indexPath: IndexPath) -> UITableViewCell {
//
//        let cell = tableView.dequeueReusableCell(withIdentifier: "cell", for: indexPath)
//        let entry = faceList[indexPath.row]
//
//        cell.textLabel?.text = entry.name
//        cell.imageView?.image = entry.image
//        cell.imageView?.contentMode = .scaleAspectFill
//        cell.imageView?.clipsToBounds = true
//        cell.imageView?.layer.cornerRadius = 8
//        cell.imageView?.layer.borderWidth = 1
//        cell.imageView?.layer.borderColor = UIColor.lightGray.cgColor
//        cell.imageView?.frame.size = CGSize(width: 60, height: 60)
//
//        return cell
//    }
//}


import UIKit

struct FaceEntry {
    let name: String
    let image: UIImage
}

final class FaceListViewController: UIViewController {

    private var faceList: [FaceEntry] = []

    private let tableView = UITableView()
    private let captureButton = UIButton(type: .system)
    private let verifyButton = UIButton(type: .system)

    override func viewDidLoad() {
        super.viewDidLoad()
        view.backgroundColor = .white
        setupUI()
    }

    private func setupUI() {
        tableView.dataSource = self
        tableView.register(FaceCell.self, forCellReuseIdentifier: "FaceCell")
        tableView.translatesAutoresizingMaskIntoConstraints = false

        captureButton.setTitle("Capture", for: .normal)
        verifyButton.setTitle("Verify", for: .normal)

        captureButton.addTarget(self, action: #selector(captureTapped), for: .touchUpInside)
        verifyButton.addTarget(self, action: #selector(verifyTapped), for: .touchUpInside)

        let buttonStack = UIStackView(arrangedSubviews: [captureButton, verifyButton])
        buttonStack.axis = .horizontal
        buttonStack.spacing = 20
        buttonStack.translatesAutoresizingMaskIntoConstraints = false

        view.addSubview(tableView)
        view.addSubview(buttonStack)

        NSLayoutConstraint.activate([
            tableView.topAnchor.constraint(equalTo: view.safeAreaLayoutGuide.topAnchor),
            tableView.leadingAnchor.constraint(equalTo: view.leadingAnchor),
            tableView.trailingAnchor.constraint(equalTo: view.trailingAnchor),

            buttonStack.topAnchor.constraint(equalTo: tableView.bottomAnchor, constant: 10),
            buttonStack.centerXAnchor.constraint(equalTo: view.centerXAnchor),
            buttonStack.bottomAnchor.constraint(equalTo: view.safeAreaLayoutGuide.bottomAnchor, constant: -10)
        ])
    }

    @objc private func captureTapped() {
        let captureVC = FaceCaptureViewController()
        captureVC.onFaceCaptured = { [weak self] name, image in
            self?.faceList.append(FaceEntry(name: name, image: image))
            DispatchQueue.main.async {
                self?.tableView.reloadData()
            }
        }
        present(captureVC, animated: true)
    }

    @objc private func verifyTapped() {
        // Optional: push verify VC
    }
}

// MARK: - UITableViewDataSource
extension FaceListViewController: UITableViewDataSource {

    func tableView(_ tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
        faceList.count
    }

    func tableView(_ tableView: UITableView, cellForRowAt indexPath: IndexPath) -> UITableViewCell {

        let cell = tableView.dequeueReusableCell(withIdentifier: "FaceCell", for: indexPath) as! FaceCell
        let entry = faceList[indexPath.row]
        cell.configure(name: entry.name, image: entry.image)
        return cell
    }
}

// MARK: - Custom TableViewCell
final class FaceCell: UITableViewCell {

    private let faceImageView: UIImageView = {
        let iv = UIImageView()
        iv.contentMode = .scaleAspectFill
        iv.clipsToBounds = true
        iv.layer.cornerRadius = 8
        iv.layer.borderWidth = 1
        iv.layer.borderColor = UIColor.lightGray.cgColor
        iv.translatesAutoresizingMaskIntoConstraints = false
        return iv
    }()

    private let nameLabel: UILabel = {
        let label = UILabel()
        label.font = .systemFont(ofSize: 16, weight: .medium)
        label.translatesAutoresizingMaskIntoConstraints = false
        return label
    }()

    override init(style: UITableViewCell.CellStyle, reuseIdentifier: String?) {
        super.init(style: style, reuseIdentifier: reuseIdentifier)

        contentView.addSubview(faceImageView)
        contentView.addSubview(nameLabel)

        NSLayoutConstraint.activate([
            faceImageView.leadingAnchor.constraint(equalTo: contentView.leadingAnchor, constant: 10),
            faceImageView.centerYAnchor.constraint(equalTo: contentView.centerYAnchor),
            faceImageView.widthAnchor.constraint(equalToConstant: 60),
            faceImageView.heightAnchor.constraint(equalToConstant: 60),

            nameLabel.leadingAnchor.constraint(equalTo: faceImageView.trailingAnchor, constant: 10),
            nameLabel.centerYAnchor.constraint(equalTo: contentView.centerYAnchor),
            nameLabel.trailingAnchor.constraint(equalTo: contentView.trailingAnchor, constant: -10)
        ])
    }

    required init?(coder: NSCoder) { fatalError() }

    func configure(name: String, image: UIImage) {
        nameLabel.text = name
        faceImageView.image = image
    }
}
